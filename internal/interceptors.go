package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// NewResponseWriter creates a wrapped http.ResponseWriter
type NewResponseWriter func(http.ResponseWriter, *http.Request) http.ResponseWriter

// Response is the structure of every http response
type Response struct {
	APIVersion string                 `json:"apiVersion"`
	ID         string                 `json:"id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Data       []interface{}          `json:"data,omitempty"`
	Errors     map[string]interface{} `json:"error,omitempty"`
}

// ResponseProperties are the predefined set of properties for each response
type ResponseProperties struct {
	APIVersion string
}

// ResponseInterceptor intercepts its http.ResponseWriter and injects Response metadata
// along with anything that is passed in the Write function
type ResponseInterceptor struct {
	ResponseProperties ResponseProperties
	Writer             http.ResponseWriter
	Request            *http.Request
}

// Header returns the original response writer headers
func (interceptor ResponseInterceptor) Header() http.Header {
	return interceptor.Writer.Header()
}

// Write writes additional response metadata to the original http.ResponseWriter
func (interceptor ResponseInterceptor) Write(b []byte) (int, error) {
	responseProperties := interceptor.ResponseProperties
	bytesWritten := 0

	ctx := interceptor.Request.Context()

	requestID, ok := ctx.Value("id").(string)
	if !ok {
		requestID = ""
	}
	method, ok := ctx.Value("method").(string)
	if !ok {
		requestID = ""
	}

	err := func() error {
		response := &Response{
			APIVersion: responseProperties.APIVersion,
			ID:         requestID,
			Method:     method,
		}
		err := json.Unmarshal(b, &response)
		if err != nil {
			return err
		}

		responseBytes := new(bytes.Buffer)
		err = json.NewEncoder(responseBytes).Encode(response)
		if err != nil {
			return err
		}

		bytesWritten, err = interceptor.Writer.Write(responseBytes.Bytes())

		return err
	}()

	return bytesWritten, err
}

// WriteHeader writes a header to the original http.ResponseWriter
func (interceptor ResponseInterceptor) WriteHeader(statusCode int) {
	interceptor.Writer.WriteHeader(statusCode)
}

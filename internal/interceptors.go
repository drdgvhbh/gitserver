package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type NewResponseWriter func(http.ResponseWriter, *http.Request) http.ResponseWriter

type Response struct {
	APIVersion string                 `json:"apiVersion"`
	ID         string                 `json:"id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Error      map[string]interface{} `json:"error,omitempty"`
}

type ResponseProperties struct {
	APIVersion string
}

type ResponseInterceptor struct {
	ResponseProperties ResponseProperties
	Writer             http.ResponseWriter
	Request            *http.Request
}

func (interceptor ResponseInterceptor) Header() http.Header {
	return interceptor.Writer.Header()
}

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

func (interceptor ResponseInterceptor) WriteHeader(statusCode int) {
	interceptor.WriteHeader(statusCode)
}

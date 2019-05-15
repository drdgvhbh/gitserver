package internal

import (
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

type Response struct {
	APIVersion string `json:"apiVersion"`
	Data map[string]interface{} `json:"data,omitempty"`
	Error map[string]interface{} `json:"error,omitempty"`
}

type ResponseInterceptor struct {
	responseProperties ResponseProperties
	writer             http.ResponseWriter
}

func (interceptor ResponseInterceptor) Header() http.Header {
	return interceptor.writer.Header()
}

func (interceptor ResponseInterceptor) Write(b []byte) (int, error) {
	responseProperties := interceptor.responseProperties
	bytesWritten := 0
	err := func() error {
		response := &Response{
			APIVersion: responseProperties.APIVersion,
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

		spew.Sdump(responseBytes.Bytes())
		bytesWritten, err = interceptor.writer.Write(responseBytes.Bytes())

		return err
	}()

	return bytesWritten, err
}

func (interceptor ResponseInterceptor) WriteHeader(statusCode int) {
	interceptor.WriteHeader(statusCode)
}
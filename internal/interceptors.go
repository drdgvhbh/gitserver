package internal

import (
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

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
		var data ResponsePayload
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}

		interceptedData := &responseProperties

		responseBytes := new(bytes.Buffer)
		err = json.NewEncoder(responseBytes).Encode(interceptedData)
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
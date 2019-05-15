package internal

import (
	"net/http"
)

type ResponseProperties struct {
	APIVersion string `json:"apiVersion"`
}

func (responseProperties ResponseProperties) Handler(handler func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		interceptedWriter := &ResponseInterceptor{responseProperties: responseProperties, writer: writer}
		handler(interceptedWriter, request)
	}
}
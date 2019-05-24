package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-git.v4"
)

var versionPrefixRegex = regexp.MustCompile("/v[0-9]+/")

// ContentTypeMiddleware injects application/json as the content type
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

// RequestMethodContextMiddleware injects the http route, along with its http method in the response
func RequestMethodContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		ctx := context.WithValue(
			request.Context(),
			"method",
			fmt.Sprintf("%s.%s",
				strings.Replace(
					versionPrefixRegex.ReplaceAllString(
						request.RequestURI, ""), "/", ".", -1),
				strings.ToLower(request.Method),
			))

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// RequestIDContextMiddleware injects a UUID as the request ID in the response
func RequestIDContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), "id", uuid.New().String())

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// NewResponseWriterMiddleware creates a middleware that intercepts another http.ResponseWriter
func NewResponseWriterMiddleware(newResponseWriter NewResponseWriter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				interceptedWriter := newResponseWriter(writer, request)

				next.ServeHTTP(interceptedWriter, request)
			})
	}
}

// OpenRepositoryMiddleware attempts to open a git repository to see if it exists,
// before passing it down the chain. If there is an error opening the repository,
// it will redirect the flow of control to the error handler
func OpenRepositoryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		repositoryPath := vars["directory"]
		repository, err := git.PlainOpen(repositoryPath)
		if err != nil {
			errorPayload := ResponsePayload{
				Error: map[string]interface{}{
					"error": err.Error(),
				},
			}
			err = json.NewEncoder(writer).Encode(&errorPayload)
			return
		}

		ctx := context.WithValue(
			request.Context(),
			"repository",
			repository)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

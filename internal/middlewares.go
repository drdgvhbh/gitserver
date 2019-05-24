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

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

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

func RequestIDContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), "id", uuid.New().String())

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func NewResponseWriterMiddleware(newResponseWriter NewResponseWriter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				interceptedWriter := newResponseWriter(writer, request)

				next.ServeHTTP(interceptedWriter, request)
			})
	}
}

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

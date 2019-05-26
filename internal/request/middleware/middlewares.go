package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/drdgvhbh/gitserver/internal/request"

	"github.com/drdgvhbh/gitserver/internal/response"

	"github.com/drdgvhbh/gitserver/internal/git"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var versionPrefixRegex = regexp.MustCompile("/v[0-9]+/")
var repositoryDoesNotExistRegex = regexp.MustCompile("repository does not exist")

func NewAuthMiddleware(apiKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			err := (func() error {
				authorization := request.Header.Get("Authorization")

				authorizationRegexp, err := regexp.Compile(
					fmt.Sprintf("^%s$", apiKey))
				if err != nil {
					return err
				}

				if authorizationRegexp.MatchString(authorization) {
					return nil
				}

				return errors.New("Unauthorized")
			})()

			if err != nil {
				errorPayload := &response.Payload{
					Errors: map[string]interface{}{
						"error": "Unauthorized",
					},
				}
				writer.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(writer).Encode(&errorPayload)
			} else {
				next.ServeHTTP(writer, request)
			}
		})
	}
}

// ContentType injects application/json as the content type
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

// MethodContext injects the http route, along with its http method in the response
func MethodContext(next http.Handler) http.Handler {
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

// IDContext injects a UUID as the request ID in the response
func IDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), "id", uuid.New().String())

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// NewResponseWriter creates a middleware that intercepts another http.ResponseWriter
func NewResponseWriter(newResponseWriter request.NewResponseWriter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				interceptedWriter := newResponseWriter(writer, request)

				next.ServeHTTP(interceptedWriter, request)
			})
	}
}

func RepositoryDirectoryVariableSanitizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		repositoryDirectory := vars["directory"]

		mux.Vars(request)["directory"] = strings.ReplaceAll(repositoryDirectory, "|", "/")

		next.ServeHTTP(writer, request)
	})
}

// NewOpenRepository creates a middleware that attempts to open a
// git repository to see if it exists, before passing it down the chain.
// If there is an error opening the repository, it will redirect the flow of
// control to the error handler
func NewOpenRepository(reader git.Reader) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			repositoryPath := vars["directory"]
			_, err := reader.Open(repositoryPath)
			if err != nil {
				var errorPayload *response.Payload
				if repositoryDoesNotExistRegex.MatchString(err.Error()) {
					errorPayload = &response.Payload{
						Errors: map[string]interface{}{
							"error": err.Error(),
						},
					}
					writer.WriteHeader(http.StatusNotFound)
				} else {
					panic(err) // TODO: Replace with a log statement or something more appropriate
				}

				err = json.NewEncoder(writer).Encode(&errorPayload)
				if err != nil {
					panic(err) // TODO: Replace with a log statement or something more appropriate
				}
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

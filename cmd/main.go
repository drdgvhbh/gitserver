//go:generate swagger generate spec -o swagger.json
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/drdgvhbh/gitserver/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var responseProperties = &internal.ResponseProperties{
	APIVersion: "0.0.1",
}

func newResponseWriter(
	writer http.ResponseWriter,
	request *http.Request,
) http.ResponseWriter {
	return &internal.ResponseInterceptor{
		ResponseProperties: *responseProperties,
		Writer:             writer,
		Request:            request,
	}
}

func main() {

	router := mux.NewRouter()
	router.Use(internal.ContentTypeMiddleware)
	router.Use(internal.RequestIDContextMiddleware)
	router.Use(internal.RequestMethodContextMiddleware)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()
	apiVersionRouter.Use(internal.NewResponseWriterMiddleware(newResponseWriter))

	repositoriesRouter := apiVersionRouter.
		PathPrefix("/repositories/{directory}").
		Subrouter()
	repositoriesRouter.Use(internal.OpenRepositoryMiddleware)

	repositoriesRouter.
		HandleFunc("/commits", internal.CommitsHandler).
		Methods("GET")

	server := &http.Server{
		Handler:      handlers.RecoveryHandler()(router),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d\n", 8000)
	log.Fatal(server.ListenAndServe())
}

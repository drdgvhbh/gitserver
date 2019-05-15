//go:generate swagger generate spec -o swagger.json
package main

import (
	"github.com/drdgvhbh/gitserver/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	responseProperties := &internal.ResponseProperties{
		APIVersion: "0.0.1",
	}

	router := mux.NewRouter()
	router.Use(internal.ContentTypeMiddleware)
	router.Use(internal.RequestContextMiddleware)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()

	repositoriesRouter := apiVersionRouter.PathPrefix("/repositories/{directory}").Subrouter()

	repositoriesRouter.HandleFunc("/commits", responseProperties.Handler(internal.CommitsHandler)).
		Methods("GET")

	server := &http.Server{
		Handler: handlers.RecoveryHandler()(router),
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

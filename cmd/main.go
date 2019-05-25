//go:generate swagger generate spec -o swagger.json
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/drdgvhbh/gitserver/internal"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

func main() {
	fs := osfs.New("")
	rootHandler := internal.NewRootHandler(fs)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d\n", 8000)
	log.Fatal(server.ListenAndServe())
}

package internal

import (
	"net/http"

	request2 "github.com/drdgvhbh/gitserver/internal/request"
	"github.com/drdgvhbh/gitserver/internal/request/middleware"

	"github.com/drdgvhbh/gitserver/internal/repository/commit"
	"github.com/drdgvhbh/gitserver/internal/response"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-billy.v4"
)

var (
	responseProperties = &response.Properties{
		APIVersion: "0.0.1",
	}
)

func newResponseWriter(
	writer http.ResponseWriter,
	request *http.Request,
) http.ResponseWriter {
	return &request2.ResponseWriter{
		ResponseProperties: *responseProperties,
		Writer:             writer,
		Request:            request,
	}
}

func NewRootHandler(worktree billy.Filesystem) http.Handler {
	fileSystem := git.NewReader(worktree)

	router := mux.NewRouter()
	router.Use(middleware.ContentType)
	router.Use(middleware.IDContext)
	router.Use(middleware.MethodContext)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()
	apiVersionRouter.Use(middleware.NewResponseWriter(newResponseWriter))

	repositoriesRouter := apiVersionRouter.
		PathPrefix("/repositories/{directory}").
		Subrouter()
	repositoriesRouter.Use(middleware.RepositoryDirectoryVariableSanitizer)
	repositoriesRouter.Use(middleware.NewOpenRepository(fileSystem))

	// swagger:route GET /repositories/{directory}/commit listCommits
	//
	// List commits
	//
	// This will list the commit in the specified repository.
	//
	//     	Consumes:
	//     	- application/json
	//
	//			Produces:
	//			- application/json
	//
	//			Schemes: http
	//
	//			Security:
	//				api_key:
	//			Responses:
	//       	200: GetCommitsOkResponse
	repositoriesRouter.
		HandleFunc("/commits", commit.NewGetCommitsHandler(fileSystem)).
		Methods("GET")

	return handlers.RecoveryHandler()(router)
}

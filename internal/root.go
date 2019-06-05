package internal

import (
	"net/http"

	"github.com/drdgvhbh/gitserver/internal/repository"

	request2 "github.com/drdgvhbh/gitserver/internal/request"
	"github.com/drdgvhbh/gitserver/internal/request/middleware"

	"github.com/drdgvhbh/gitserver/internal/repository/commit"
	"github.com/drdgvhbh/gitserver/internal/repository/reference"
	"github.com/drdgvhbh/gitserver/internal/response"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-billy.v4"
)

var _ = repository.Params{}

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
	apiVersionRouter.Use(middleware.NewAuthMiddleware("e8b8dc29-d1d9-495d-b509-4dde3701018b"))

	repositoriesRouter := apiVersionRouter.
		PathPrefix("/repositories/{directory}").
		Subrouter()
	repositoriesRouter.Use(middleware.RepositoryDirectoryVariableSanitizer)
	repositoriesRouter.Use(middleware.NewOpenRepository(fileSystem))

	// swagger:route GET /repositories/{directory}/commits listCommits
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

	commitRouter := repositoriesRouter.
		PathPrefix("/commits/{hash}").
		Subrouter()

	// swagger:route GET /repositories/{directory}/commit/{hash} getCommit
	//
	// Get commit
	//
	// This will get the specified commit in the specified repository.
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
	//       	200: GetCommitOkResponse
	commitRouter.
		HandleFunc("", commit.NewGetCommitHandler(fileSystem)).
		Methods("GET")
	commitRouter.
		HandleFunc("/diff", commit.NewGetCommitDiffHandler(fileSystem)).
		Methods("GET")

	// swagger:route GET /repositories/{directory}/references listReferences
	//
	// List references
	//
	// This will list the references in the specified repository.
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
	//       	200: GetReferencesOkResponse
	repositoriesRouter.
		HandleFunc("/references", reference.NewGetReferencesHandler(fileSystem)).
		Methods("GET")

	return handlers.RecoveryHandler()(router)
}

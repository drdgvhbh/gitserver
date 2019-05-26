package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4"
)

// ResponsePayload is the payload for every response
type ResponsePayload struct {
	Data  []interface{}          `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type CommitContributor struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type CommitLogData struct {
	Hash    string             `json:"hash,omitempty"`
	Summary string             `json:"summary,omitempty"`
	Author  *CommitContributor `json:"author,omitempty"`
}

var (
	responseProperties = &ResponseProperties{
		APIVersion: "0.0.1",
	}
)

func newResponseWriter(
	writer http.ResponseWriter,
	request *http.Request,
) http.ResponseWriter {
	return &ResponseInterceptor{
		ResponseProperties: *responseProperties,
		Writer:             writer,
		Request:            request,
	}
}

func NewRootHandler(worktree billy.Filesystem) http.Handler {
	fileSystem := git.NewReader(worktree)

	router := mux.NewRouter()
	router.Use(ContentTypeMiddleware)
	router.Use(RequestIDContextMiddleware)
	router.Use(RequestMethodContextMiddleware)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()
	apiVersionRouter.Use(NewResponseWriterMiddleware(newResponseWriter))

	repositoriesRouter := apiVersionRouter.
		PathPrefix("/repositories/{directory}").
		Subrouter()
	repositoriesRouter.Use(RepositoryDirectoryVariableSanitizer)
	repositoriesRouter.Use(NewOpenRepositoryMiddleware(fileSystem))

	repositoriesRouter.
		HandleFunc("/commits", NewCommitsHandler(fileSystem)).
		Methods("GET")

	return handlers.RecoveryHandler()(router)
}

// CommitsHandler returns the git commits in the specified repository
func NewCommitsHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		repositoryPath := vars["directory"]
		logrus.Println(vars["directory"])
		repository, _ := reader.Open(repositoryPath)

		err := (func() error {
			head, err := repository.Head()
			if err != nil {
				return err
			}

			commitHistory, err := repository.Log(
				&git.LogOptions{From: head.Hash()})
			if err != nil {
				return err
			}

			var commitLogData []interface{}

			err = commitHistory.ForEach(func(commit git.Commit) error {
				commitLogData = append(commitLogData, CommitLogData{
					Hash:    commit.Hash(),
					Summary: commit.Summary(),
					Author: &CommitContributor{
						Name:      commit.Author().Name(),
						Email:     commit.Author().Email(),
						Timestamp: commit.Author().Timestamp().Format(time.RFC3339),
					},
				})
				return nil
			})
			if err != nil {
				return err
			}

			dataPayload := ResponsePayload{
				Data: commitLogData,
			}

			return json.NewEncoder(writer).Encode(&dataPayload)
		})()

		if err != nil {
			errorPayload := ResponsePayload{
				Error: map[string]interface{}{
					"error": err,
				},
			}
			err = json.NewEncoder(writer).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
		}
	}
}

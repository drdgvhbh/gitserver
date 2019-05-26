package commit

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/gorilla/mux"
)

// List of commits in the repository
// swagger:response GetCommitsOkResponse
type GetCommitsOKResponse struct {
	// in: body
	Body struct {
		response.Base
		Data []LogData
	}
}

// CommitsHandler returns the git commit in the specified repository
func NewGetCommitsHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		repositoryPath := vars["directory"]
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
				author := commit.Author()
				committer := commit.Committer()

				commitLogData = append(commitLogData, LogData{
					Hash:    commit.Hash(),
					Summary: commit.Summary(),
					Author: &Contributor{
						Name:      author.Name(),
						Email:     author.Email(),
						Timestamp: author.Timestamp().Format(time.RFC3339),
					},
					Committer: &Contributor{
						Name:      committer.Name(),
						Email:     committer.Email(),
						Timestamp: committer.Timestamp().Format(time.RFC3339),
					},
				})
				return nil
			})
			if err != nil {
				return err
			}

			dataPayload := response.Payload{
				Data: commitLogData,
			}

			return json.NewEncoder(writer).Encode(&dataPayload)
		})()

		if err != nil {
			errorPayload := response.Payload{
				Errors: map[string]interface{}{
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

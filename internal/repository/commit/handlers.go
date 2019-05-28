package commit

import (
	"encoding/json"
	"net/http"
	"sort"
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
		// The request method
		//
		// required: true
		// example: repositories.%7Chome%7Cdrd%7Cgo%7Csrc%7Cgithub.com%7Cdrdgvhbh%7Cgitserver.commits.get
		Method string `json:"method,omitempty"`
		// The response data
		//
		// required: true
		Data []LogData `json:"data,omitempty"`
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

			var commitLogData []LogData

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

			sort.Slice(commitLogData, func(i, j int) bool {
				leftTimestamp, _ := time.Parse(time.RFC3339, commitLogData[i].Committer.Timestamp)
				rightTimestamp, _ := time.Parse(time.RFC3339, commitLogData[j].Committer.Timestamp)
				diff := leftTimestamp.Sub(rightTimestamp)

				return diff.Seconds() > 0
			})

			data := make([]interface{}, len(commitLogData))
			for i := range commitLogData {
				data[i] = commitLogData[i]
			}

			dataPayload := response.Payload{
				Data: data,
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

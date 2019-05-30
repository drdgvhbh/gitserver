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
		Data []Commit `json:"data,omitempty"`
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

			references := make(map[string][]string)

			refIter, err := repository.References()
			if err != nil {
				return err
			}

			defer refIter.Close()

			_ = refIter.ForEach(func(ref git.Reference) error {
				hash := ref.Hash().String()
				references[hash] = append(references[hash], string(ref.Name()))

				return nil
			})

			var commitData []Commit

			defer commitHistory.Close()
			err = commitHistory.ForEach(func(commit git.Commit) error {
				author := commit.Author()
				committer := commit.Committer()

				commitData = append(commitData, Commit{
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
					References: references[commit.Hash()],
				})
				return nil
			})
			if err != nil {
				return err
			}

			sort.Slice(commitData, func(i, j int) bool {
				leftTimestamp, _ := time.Parse(time.RFC3339, commitData[i].Committer.Timestamp)
				rightTimestamp, _ := time.Parse(time.RFC3339, commitData[j].Committer.Timestamp)
				diff := leftTimestamp.Sub(rightTimestamp)

				return diff.Seconds() > 0
			})

			data := make([]interface{}, len(commitData))
			for i := range commitData {
				data[i] = commitData[i]
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

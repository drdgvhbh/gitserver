package commit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
				commitRefs := references[commit.Hash()]
				commitData = append(commitData, *NewCommit(commit, commitRefs))
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

func NewGetCommitHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		repositoryPath := vars["directory"]
		commitHash := vars["hash"]
		repository, _ := reader.Open(repositoryPath)

		opts := &git.LogOptions{
			From: git.NewHash(commitHash),
		}

		err := (func() error {
			commitHistory, err := repository.Log(opts)
			if err != nil {
				return err
			}

			specifiedCommit, err := commitHistory.Next()
			if err != nil {
				logrus.Println(err)
				return err
			}
			if specifiedCommit == nil {
				w.WriteHeader(http.StatusNotFound)
				return fmt.Errorf("commit '%s' not found", commitHash)
			}

			refs, err := repository.ReferenceMap()
			if err != nil {
				return err
			}
			refNames := refs[specifiedCommit.Hash()].MapStr(func(ref git.Reference) string {
				return string(ref.Name())
			})

			var commitData []Commit
			commitData = append(commitData, *NewCommit(specifiedCommit, refNames))

			data := make([]interface{}, len(commitData))
			for i := range commitData {
				data[i] = commitData[i]
			}

			dataPayload := response.Payload{
				Data: data,
			}

			_ = json.NewEncoder(w).Encode(&dataPayload)

			return nil
		})()

		if err != nil {
			errPayload := response.Payload{
				Errors: map[string]interface{}{
					"error": err.Error(),
				},
			}
			_ = json.NewEncoder(w).Encode(&errPayload)
		}
	}
}

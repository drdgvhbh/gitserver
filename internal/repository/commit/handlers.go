package commit

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

		var err error
		defer (func() {
			writeError(err, writer)
		})()

		err = (func() error {
			head, err := repository.Head()
			if err != nil {
				return err
			}

			commitHistory, err := repository.Log(
				&git.LogOptions{From: head.Hash()})
			if err != nil {
				return err
			}

			hashToRefNames, err := hashToRefNames(repository)
			if err != nil {
				return err
			}

			var commitData []Commit

			defer commitHistory.Close()
			err = commitHistory.ForEach(func(commit git.Commit) error {
				commitRefs := hashToRefNames[commit.Hash()]
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

			writeData(data, writer)

			return nil
		})()
	}
}

// Gets the specified commit in the repository
// swagger:response GetCommitOkResponse
type GetCommitOKResponse struct {
	// in: body
	Body struct {
		response.Base
		// The request method
		//
		// required: true
		// example: repositories.%7Chome%7Cdrd%7Cgo%7Csrc%7Cgithub.com%7Cdrdgvhbh%7Cgitserver.commits.5dd5708f4c284919b1ef22f44e5d98f7d7579910.get
		Method string `json:"method,omitempty"`
		// The response data
		//
		// required: true
		Data []Commit `json:"data,omitempty"`
	}
}

func NewGetCommitHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		repositoryPath := vars["directory"]
		commitHash := vars["hash"]
		repository, _ := reader.Open(repositoryPath)

		var err error
		defer (func() {
			writeError(err, w)
		})()

		err = (func() error {
			specifiedCommit, err := repository.FindCommit(git.NewHash(commitHash))

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return fmt.Errorf("commit '%s' not found", commitHash)
			}

			hashToRefNames, err := hashToRefNames(repository)
			if err != nil {
				return err
			}
			commitRefNames := hashToRefNames[specifiedCommit.Hash()]

			var commitData []Commit
			commitData = append(commitData, *NewCommit(specifiedCommit, commitRefNames))

			data := make([]interface{}, len(commitData))
			for i := range commitData {
				data[i] = commitData[i]
			}

			writeData(data, w)

			return nil
		})()
	}
}

func NewGetCommitDiffHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		repositoryPath := vars["directory"]
		commitHash := vars["hash"]
		repository, _ := reader.Open(repositoryPath)

		var err error
		defer (func() {
			writeError(err, w)
		})()

		err = (func() error {
			changes, err := repository.Diff(git.NewHash(commitHash))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return errors.Wrapf(err, "commit '%s' not found", commitHash)
			}

			var changeData []*Change

			for _, change := range changes {
				changeType, err := change.Action()
				if err != nil {
					return errors.Wrapf(err, "change has no type")
				}
				changeTypeStr, err := changeType.String()
				if err != nil {
					return errors.Wrapf(err, "change cannot be converted to string")
				}

				var path string
				switch changeType {
				case git.Insert:
					path = change.FilePathAfter()
				case git.Modify:
					path = change.FilePathAfter()
				case git.Delete:
					path = change.FilePathBefore()
				}

				changeData = append(changeData, NewChange(changeTypeStr, path))
			}

			data := make([]interface{}, len(changeData))
			for i := range changeData {
				data[i] = changeData[i]
			}
			writeData(data, w)

			return nil
		})()

	}
}

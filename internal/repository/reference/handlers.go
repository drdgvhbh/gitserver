package reference

import (
	"encoding/json"
	"net/http"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/gorilla/mux"
)

// List of references in the repository
// swagger:response GetReferencesOkResponse
type GetReferencesOKResponse struct {
	// in: body
	Body struct {
		response.Base
		// The request method
		//
		// required: true
		// example: repositories.%7Chome%7Cdrd%7Cgo%7Csrc%7Cgithub.com%7Cdrdgvhbh%7Cgitserver.references.get
		Method string `json:"method,omitempty"`
		// The response data
		//
		// required: true
		Data []Reference `json:"data,omitempty"`
	}
}

func NewGetReferencesHandler(reader git.Reader) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		repositoryPath := vars["directory"]
		repository, _ := reader.Open(repositoryPath)

		err := (func() error {
			referIter, err := repository.References()
			if err != nil {
				return err
			}

			var referencesData []Reference

			_ = referIter.ForEach(func(ref git.Reference) error {
				referencesData = append(referencesData, Reference{
					Hash: ref.Hash().String(),
					Name: string(ref.Name()),
				})

				return nil
			})

			data := make([]interface{}, len(referencesData))
			for i := range referencesData {
				data[i] = referencesData[i]
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

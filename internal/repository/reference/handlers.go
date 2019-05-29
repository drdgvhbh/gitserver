package reference

import (
	"encoding/json"
	"net/http"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/gorilla/mux"
)

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

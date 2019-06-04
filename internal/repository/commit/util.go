package commit

import (
	"encoding/json"
	"net/http"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/sirupsen/logrus"
)

func hashToRefNames(r git.Repository) (map[string][]string, error) {
	refMap := make(map[string][]string)

	refIter, err := r.References()
	if err != nil {
		return nil, err
	}

	defer refIter.Close()

	_ = refIter.ForEach(func(ref git.Reference) error {
		hash := ref.Hash().String()
		refMap[hash] = append(refMap[hash], string(ref.Name()))

		return nil
	})

	return refMap, nil
}

func writeError(err error, w http.ResponseWriter) {
	if err != nil {
		errorPayload := response.Payload{
			Errors: map[string]interface{}{
				"error": err,
			},
		}
		err = json.NewEncoder(w).Encode(&errorPayload)
		if err != nil {
			logrus.Warnln(err)
		}
	}
}

package commit

import (
	"github.com/drdgvhbh/gitserver/internal/git"
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

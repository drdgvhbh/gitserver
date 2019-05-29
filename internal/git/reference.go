package git

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type ReferenceName string

type Reference interface {
	Hash() Hash
	Name() ReferenceName
}

type ReferenceIter interface {
	Next() (Reference, error)
	ForEach(func(Reference) error) error
	Close()
}

type GitReferenceIter struct {
	Wrapee storer.ReferenceIter
}

func (iter *GitReferenceIter) Next() (Reference, error) {
	reference, err := iter.Wrapee.Next()

	if err != nil {
		return nil, err
	}

	return &GitReference{
		Wrapee: reference,
	}, nil
}

func (iter *GitReferenceIter) ForEach(fn func(Reference) error) error {
	return iter.Wrapee.ForEach(func(reference *plumbing.Reference) error {
		return fn(&GitReference{
			Wrapee: reference,
		})
	})
}

func (iter *GitReferenceIter) Close() {
	iter.Wrapee.Close()
}

type GitReference struct {
	Wrapee *plumbing.Reference
}

func (ref *GitReference) Hash() Hash {
	return Hash(ref.Wrapee.Hash())
}

func (ref *GitReference) Name() ReferenceName {
	return ReferenceName(ref.Wrapee.Name())
}


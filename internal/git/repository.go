package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Hash [20]byte

func (hash Hash) String() string {
	return plumbing.Hash(hash).String()
}

func (hash Hash) IsZero() bool {
	return plumbing.Hash(hash).IsZero()
}

type LogOptions struct {
	From Hash
}

type Repository interface {
	Head() (Reference, error)
	Log(options *LogOptions) (CommitIter, error)
	Reference(name ReferenceName) (Reference, error)
	References() (ReferenceIter, error)
}

type GitRepository struct {
	Wrapee *git.Repository
}

func (repo *GitRepository) Head() (Reference, error) {
	head, err := repo.Wrapee.Head()

	if err != nil {
		return nil, err
	}

	return &GitReference{Wrapee: head}, nil
}

func (repo *GitRepository) Log(options *LogOptions) (CommitIter, error) {
	wrappedOpts := &git.LogOptions{
		From: plumbing.Hash(options.From),
	}

	commitIter, err := repo.Wrapee.Log(wrappedOpts)
	if err != nil {
		return nil, err
	}

	return &GitCommitIter{Wrapee: commitIter}, nil
}

func (repo *GitRepository) References() (ReferenceIter, error) {
	wrappedIter, err := repo.Wrapee.References()

	if err != nil {
		return nil, err
	}

	return &GitReferenceIter{Wrapee: wrappedIter}, nil

}

func (repo *GitRepository) Reference(name ReferenceName) (Reference, error) {
	wrapped, err := repo.Wrapee.Reference(plumbing.ReferenceName(name), false)

	if err != nil {
		return nil, err
	}

	return &GitReference{Wrapee: wrapped}, nil
}

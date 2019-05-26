package git

import (
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Hash [20]byte

type Signature interface {
	Name() string
	Email() string
	Timestamp() time.Time
}

type SignatureWrapper struct {
	Wrapee object.Signature
}

func (s SignatureWrapper) Name() string {
	return s.Wrapee.Name
}

func (s SignatureWrapper) Email() string {
	return s.Wrapee.Email
}

func (s SignatureWrapper) Timestamp() time.Time {
	return s.Wrapee.When
}

type LogOptions struct {
	From Hash
}

type Reference interface {
	Hash() Hash
}

type Commit interface {
	Summary() string
	Hash() string
	Author() Signature
}

type GitCommit struct {
	Wrapee *object.Commit
}

func (commit *GitCommit) Hash() string {
	return commit.Wrapee.Hash.String()
}

func (commit *GitCommit) Summary() string {
	message := strings.Split(commit.Wrapee.Message, "\n")

	if len(message) <= 0 {
		return ""
	}

	return message[0]
}

func (commit *GitCommit) Author() Signature {
	return &SignatureWrapper{commit.Wrapee.Author}
}

type CommitIter interface {
	Next() (Commit, error)
	ForEach(func(Commit) error) error
	Close()
}

type GitCommitIter struct {
	Wrapee object.CommitIter
}

func (iter *GitCommitIter) Next() (Commit, error) {
	commit, err := iter.Wrapee.Next()

	if err != nil {
		return nil, err
	}

	return &GitCommit{
		Wrapee: commit,
	}, nil
}

func (iter *GitCommitIter) ForEach(fn func(Commit) error) error {
	return iter.Wrapee.ForEach(func(commit *object.Commit) error {
		return fn(&GitCommit{
			Wrapee: commit,
		})
	})
}

func (iter *GitCommitIter) Close() {
	iter.Wrapee.Close()
}

type Repository interface {
	Head() (Reference, error)
	Log(options *LogOptions) (CommitIter, error)
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

type GitReference struct {
	Wrapee *plumbing.Reference
}

func (ref *GitReference) Hash() Hash {
	return Hash(ref.Wrapee.Hash())
}

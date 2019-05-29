package git

import (
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"strings"
)

type Commit interface {
	Summary() string
	Hash() string
	Author() Signature
	Committer() Signature
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

func (commit *GitCommit) Committer() Signature {
	return &SignatureWrapper{commit.Wrapee.Committer}
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


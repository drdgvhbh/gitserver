package commit

import (
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
)

type Contributor struct {
	// Contributor's Name
	//
	// required: true
	// example: Ryan Lee
	Name string `json:"name,omitempty"`
	// Contributor's email
	//
	// required: true
	// example: ryanleecode@gmail.com
	Email string `json:"email,omitempty"`
	// Timestamp of contribution
	//
	// required: true
	// example: 2019-05-26T12:41:18-04:00
	Timestamp string `json:"timestamp,omitempty"`
}

type Commit struct {
	// The hash of the commit
	//
	// required: true
	// example: e38e2cde1fada4a738f2461b283e561bc767568b
	Hash string `json:"hash,omitempty"`

	// The summary of the commit
	//
	// example: Deletes swagger documentation from the repository
	Summary string `json:"summary,omitempty"`

	// The author of the commit
	//
	// required: true
	Author *Contributor `json:"author,omitempty"`

	// The committer of the commit
	//
	// required: true
	Committer *Contributor `json:"committer,omitempty"`

	// The references pointing to this commit
	//
	// required: true
	References []string `json:"references"`
}

func NewCommit(c git.Commit, ref []string) *Commit {
	commitRefs := ref
	if commitRefs == nil {
		commitRefs = make([]string, 0)
	}

	author := c.Author()
	committer := c.Committer()

	return &Commit{
		Hash:    c.Hash(),
		Summary: c.Summary(),
		Author: &Contributor{
			Name:      author.Name(),
			Email:     author.Email(),
			Timestamp: author.Timestamp().Format(time.RFC3339),
		},
		Committer: &Contributor{
			Name:      committer.Name(),
			Email:     committer.Email(),
			Timestamp: committer.Timestamp().Format(time.RFC3339),
		},
		References: commitRefs,
	}
}

func NewChange(changeType string, path string) *Change {
	return &Change{
		Type: changeType,
		Path: path,
	}
}

type Change struct {
	// The type of change
	//
	// enum: MODIFY,INSERT,DELETE
	// required: true
	Type string `json:"type"`
	// The relative file path in the repository
	//
	// required: true
	// example: lib/dart/example.dart
	Path string `json:"path"`
}

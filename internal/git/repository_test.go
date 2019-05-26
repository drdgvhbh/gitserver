package git_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestCommitSummary(t *testing.T) {
	assert := assert.New(t)

	summary := "This is a summary"
	depCommit := object.Commit{
		Message: fmt.Sprintf("%s\n\nThis is a description", summary),
	}

	commit := git.GitCommit{
		Wrapee: &depCommit,
	}

	assert.EqualValues(summary, commit.Summary())
}

func TestCommitHash(t *testing.T) {
	assert := assert.New(t)

	hash := "55245d63089b55144010e408895a4350e163b49e"
	depCommit := object.Commit{
		Hash: plumbing.NewHash(hash),
	}

	commit := git.GitCommit{
		Wrapee: &depCommit,
	}

	assert.EqualValues(hash, commit.Hash())
}

func TestSignatureWrapperName(t *testing.T) {
	assert := assert.New(t)

	name := "Ryan Lee"
	depSignature := object.Signature{
		Name: name,
	}

	signature := git.SignatureWrapper{
		Wrapee: depSignature,
	}

	assert.EqualValues(name, signature.Name())
}

func TestSignatureWrapperEmail(t *testing.T) {
	assert := assert.New(t)

	email := "drdgvhbh@gmail.com"
	depEmail := object.Signature{
		Email: email,
	}

	signature := git.SignatureWrapper{
		Wrapee: depEmail,
	}

	assert.EqualValues(email, signature.Email())
}

func TestSignatureWrapperTimestamp(t *testing.T) {
	assert := assert.New(t)

	timestamp, err := time.Parse(time.RFC3339, "2019-05-25T15:18:47-04:00")
	assert.NoError(err)
	depTimestamp := object.Signature{
		When: timestamp,
	}

	signature := git.SignatureWrapper{
		Wrapee: depTimestamp,
	}

	assert.EqualValues(timestamp, signature.Timestamp())
}

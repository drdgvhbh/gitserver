package mock

import (
	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/stretchr/testify/mock"
)

type CommitIter struct {
	mock.Mock
}

func (c *CommitIter) Next() (git.Commit, error) {
	args := c.Called()

	return nil, args.Error(1)
}

func (c *CommitIter) ForEach(fn func(git.Commit) error) error {
	args := c.Called(fn)

	return args.Error(0)
}

func (c *CommitIter) Close() {
}

type Reference struct {
	mock.Mock
}

func (*Reference) Hash() git.Hash {
	return [20]byte{}
}

type Repository struct {
	mock.Mock
}

func (r *Repository) Head() (git.Reference, error) {
	args := r.Called()

	return &Reference{}, args.Error(1)
}

func (r *Repository) Log(options *git.LogOptions) (git.CommitIter, error) {
	args := r.Called(options)

	return &CommitIter{}, args.Error(1)
}

type Reader struct {
	mock.Mock
}

func (m *Reader) Open(path string) (git.Repository, error) {
	args := m.Called(path)

	return &Repository{}, args.Error(1)
}

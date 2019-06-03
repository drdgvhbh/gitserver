package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type File interface {
	Name() string
}

type GitFile struct {
	w *object.File
}

func (f *GitFile) Name() string {
	return f.w.Name
}

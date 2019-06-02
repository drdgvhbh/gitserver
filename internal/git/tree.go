package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type Tree interface{}

type GitTree struct {
	Wrapee *object.Tree
}

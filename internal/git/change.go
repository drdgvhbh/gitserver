package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type Changes = []Change

type Change interface{}

type GitChange struct {
	Wrapee *object.Change
}

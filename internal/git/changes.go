package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type Changes interface{}

type GitChanges struct {
	Wrapee object.Changes
}

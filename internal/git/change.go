package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type Action int

const (
	_ Action = iota
	Insert
	Delete
	Modify
)

type Changes = []Change

type Change interface{}

type GitChange struct {
	Wrapee *object.Change
}

func (c *GitChange) Action() (Action, error) {
	action, err := c.Wrapee.Action()

	if err != nil {
		return 0, err
	}

	return Action(action), nil
}

func (c *GitChange) Files() (from, to File, err error) {
	f, t, err := c.Wrapee.Files()

	return &GitFile{w: f}, &GitFile{w: t}, err
}

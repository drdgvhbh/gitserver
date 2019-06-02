package git

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type Tree interface {
	Diff(Tree) (Changes, error)
}

type GitTree struct {
	Wrapee *object.Tree
}

func (tree *GitTree) Diff(otherTree Tree) (Changes, error) {
	changes, err := tree.Wrapee.Diff(otherTree.(*GitTree).Wrapee)

	if err != nil {
		return nil, err
	}

	return &GitChanges{Wrapee: changes}, nil
}

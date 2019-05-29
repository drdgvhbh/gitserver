package git

import (
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

type Signature interface {
	Name() string
	Email() string
	Timestamp() time.Time
}

type SignatureWrapper struct {
	Wrapee object.Signature
}

func (s SignatureWrapper) Name() string {
	return s.Wrapee.Name
}

func (s SignatureWrapper) Email() string {
	return s.Wrapee.Email
}

func (s SignatureWrapper) Timestamp() time.Time {
	return s.Wrapee.When
}


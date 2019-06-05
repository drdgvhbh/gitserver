package git

import (
	"fmt"
)

type Action int

const (
	_ Action = iota
	Insert
	Delete
	Modify
)

func (a Action) String() (string, error) {
	switch a {
	case Insert:
		return "INSERT", nil
	case Modify:
		return "MODIFY", nil
	case Delete:
		return "DELETE", nil
	default:
		return "", fmt.Errorf("%d is not an action", a)
	}
}

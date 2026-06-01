package task

import (
	"errors"
)

var (
	ErrMismatchedOwner = errors.New("owner is not the original owner")
	ErrUnsupported     = errors.New("unsupported")
)

package task

import (
	"fmt"
	"strings"
)

type Status int32

const (
	StatusUnknown Status = iota
	StatusPending
	StatusCompleted
)

func (s Status) Valid() bool {
	return StatusUnknown < s && s <= StatusCompleted
}

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

func ParseStatus(status string) (Status, error) {
	switch strings.ToLower(status) {
	case "pending":
		return StatusPending, nil
	case "completed":
		return StatusCompleted, nil
	default:
		return StatusUnknown, fmt.Errorf("%w: %s", ErrInvalidData, status)
	}
}

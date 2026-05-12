package task

import (
	"errors"
	"fmt"
	"strings"
)

var InvalidStatus = errors.New("invalid task status")

type Status int8

const (
	StatusUnknown Status = iota
	StatusPending
	StatusCompleted
)

func (s Status) Valid() bool {
	return s != StatusUnknown
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
		return StatusUnknown, fmt.Errorf("%w: %s", InvalidStatus, status)
	}
}

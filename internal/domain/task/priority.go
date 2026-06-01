package task

import (
	"fmt"
	"strings"
)

type Priority int32

const (
	PriorityUnknown Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

func (p Priority) Valid() bool {
	return PriorityUnknown < p && p <= PriorityHigh
}

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	default:
		return "unknown"
	}
}

func ParsePriority(priority string) (Priority, error) {
	switch strings.ToLower(priority) {
	case "low":
		return PriorityLow, nil
	case "medium":
		return PriorityMedium, nil
	case "high":
		return PriorityHigh, nil
	default:
		return PriorityUnknown, fmt.Errorf("%w: %s", ErrInvalidData, priority)
	}
}

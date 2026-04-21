package model

import (
	"time"
)

type TaskStatus int

const (
	TASK_STATUS_UNKNOWN TaskStatus = iota
	TASK_STATUS_ACTIVE
	TASK_STATUS_COMPLETED
	TASK_STATUS_PENDING_DELETION
)

func (ts TaskStatus) String() string {
	switch ts {
	case TASK_STATUS_ACTIVE:
		return "active"
	case TASK_STATUS_COMPLETED:
		return "completed"
	case TASK_STATUS_PENDING_DELETION:
		return "deleting"
	default:
		return "unknown"
	}
}

type TaskPriority int

const (
	TASK_PRIORITY_UNKNOWN TaskPriority = iota
	TASK_PRIORITY_LOW
	TASK_PRIORITY_MEDIUM
	TASK_PRIORITY_HIGH
)

func (tp TaskPriority) String() string {
	switch tp {
	case TASK_PRIORITY_LOW:
		return "low"
	case TASK_PRIORITY_MEDIUM:
		return "medium"
	case TASK_PRIORITY_HIGH:
		return "high"
	default:
		return "unknown"
	}
}

type GeoPoint struct {
	Latitude, Longitude float64
}

type Task struct {
	Title       string       `validate:"required,max=200"`
	Description string       `validate:"required,max=10000"`
	Location    GeoPoint     `validate:"omitempty"`
	Priority    TaskPriority `validate:"required"`
	Status      TaskStatus   `validate:"required"`
	IsFavorite  bool         `validate:"omitempty"`
	CreatedAt   time.Time    `validate:"required"`
	UpdatedAt   time.Time    `validate:"omitempty"`
	CompletedAt time.Time    `validate:"omitempty"`
	Deadline    time.Time    `validate:"omitempty"`
}

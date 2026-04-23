package model

import (
	"time"
)

type TaskStatus int

const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusActive
	TaskStatusCompleted
	TaskStatusPendingDeletion
)

func (ts TaskStatus) String() string {
	switch ts {
	case TaskStatusActive:
		return "active"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusPendingDeletion:
		return "deleting"
	default:
		return "unknown"
	}
}

type TaskPriority int

const (
	TaskPriorityUnknown TaskPriority = iota
	TaskPriorityLow
	TaskPriorityMedium
	TaskPriorityHigh
)

func (tp TaskPriority) String() string {
	switch tp {
	case TaskPriorityLow:
		return "low"
	case TaskPriorityMedium:
		return "medium"
	case TaskPriorityHigh:
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

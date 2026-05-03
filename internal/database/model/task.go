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
	Title       string       `validate:"required,max=200" bson:"title"`
	Description string       `validate:"required,max=10000" bson:"description"`
	Location    GeoPoint     `validate:"omitempty" bson:"location"`
	Priority    TaskPriority `validate:"required" bson:"priority"`
	Status      TaskStatus   `validate:"required" bson:"status"`
	IsFavorite  bool         `validate:"omitempty" bson:"is_favorite,omitempty"`
	CreatedAt   time.Time    `validate:"required" bson:"created_at"`
	UpdatedAt   time.Time    `validate:"omitempty" bson:"updated_at,omitempty"`
	CompletedAt time.Time    `validate:"omitempty" bson:"completed_at,omitempty"`
	Deadline    time.Time    `validate:"omitempty" bson:"deadline,omitempty"`
}

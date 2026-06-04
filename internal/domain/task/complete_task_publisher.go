package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CompleteTaskEvent struct {
	TaskID      uuid.UUID
	UserID      uuid.UUID
	CompletedAt time.Time
}

func NewCompleteTaskEvent(taskID, userID uuid.UUID, completedAt time.Time) *CompleteTaskEvent {
	return &CompleteTaskEvent{
		TaskID:      taskID,
		UserID:      userID,
		CompletedAt: completedAt,
	}
}

type CompleteTaskPublisher interface {
	Publish(ctx context.Context, event *CompleteTaskEvent) error
}

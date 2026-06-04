package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type NotificationEvent struct {
	TaskID   uuid.UUID
	UserID   uuid.UUID
	Title    string
	NotifyAt []time.Time
}

func NewNotificationEvent(taskID, userID uuid.UUID, title string, notifyAt []time.Time) *NotificationEvent {
	return &NotificationEvent{
		TaskID:   taskID,
		UserID:   userID,
		Title:    title,
		NotifyAt: notifyAt,
	}
}

type NotificationEventPublisher interface {
	Notify(ctx context.Context, event *NotificationEvent) error
}

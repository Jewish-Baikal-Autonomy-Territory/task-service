package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type CreatedTaskEvent struct {
	TaskID    uuid.UUID
	UserID    uuid.UUID
	GroupID   mo.Option[uuid.UUID]
	Location  mo.Option[GeoPoint]
	CreatedAt time.Time
}

type CreatedTaskEventBuilder struct {
	taskID    uuid.UUID
	userID    uuid.UUID
	groupID   mo.Option[uuid.UUID]
	location  mo.Option[GeoPoint]
	createdAt time.Time
}

func (b *CreatedTaskEventBuilder) WithTaskID(id uuid.UUID) *CreatedTaskEventBuilder {
	b.taskID = id
	return b
}

func (b *CreatedTaskEventBuilder) WithUserID(id uuid.UUID) *CreatedTaskEventBuilder {
	b.userID = id
	return b
}

func (b *CreatedTaskEventBuilder) WithGroupID(id uuid.UUID) *CreatedTaskEventBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *CreatedTaskEventBuilder) WithLocation(location GeoPoint) *CreatedTaskEventBuilder {
	b.location = mo.Some(location)
	return b
}

func (b *CreatedTaskEventBuilder) WithCreatedAt(createdAt time.Time) *CreatedTaskEventBuilder {
	b.createdAt = createdAt
	return b
}

func (b *CreatedTaskEventBuilder) Build() *CreatedTaskEvent {
	return &CreatedTaskEvent{
		TaskID:    b.taskID,
		UserID:    b.userID,
		GroupID:   b.groupID,
		Location:  b.location,
		CreatedAt: b.createdAt,
	}
}

func NewCreatedTaskEventBuilder() *CreatedTaskEventBuilder {
	return &CreatedTaskEventBuilder{}
}

type CreatedTaskPublisher interface {
	Publish(ctx context.Context, event *CreatedTaskEvent) error
}

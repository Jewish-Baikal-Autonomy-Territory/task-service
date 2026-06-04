package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type CreateTaskCommand struct {
	OwnerID uuid.UUID
	GroupID mo.Option[uuid.UUID]

	Title       string
	Description string
	Location    mo.Option[domaintask.GeoPoint]
	IsFavorite  bool
	Priority    domaintask.Priority
	Icon        domaintask.Icon
	Deadline    mo.Option[time.Time]
	NotifyAt    []time.Time
}

type CreateTaskCommandBuilder struct {
	ownerID uuid.UUID
	groupID mo.Option[uuid.UUID]

	location    mo.Option[domaintask.GeoPoint]
	title       string
	description string
	isFavorite  bool
	priority    domaintask.Priority
	icon        domaintask.Icon
	deadline    mo.Option[time.Time]
	notifyAt    []time.Time
}

func NewCreateTaskCommandBuilder() *CreateTaskCommandBuilder {
	return &CreateTaskCommandBuilder{}
}

func (b *CreateTaskCommandBuilder) WithOwnerID(id uuid.UUID) *CreateTaskCommandBuilder {
	b.ownerID = id
	return b
}

func (b *CreateTaskCommandBuilder) WithGroupID(id uuid.UUID) *CreateTaskCommandBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *CreateTaskCommandBuilder) WithTitle(title string) *CreateTaskCommandBuilder {
	b.title = title
	return b
}

func (b *CreateTaskCommandBuilder) WithDescription(description string) *CreateTaskCommandBuilder {
	b.description = description
	return b
}

func (b *CreateTaskCommandBuilder) WithLocation(location domaintask.GeoPoint) *CreateTaskCommandBuilder {
	b.location = mo.Some(location)
	return b
}

func (b *CreateTaskCommandBuilder) WithIsFavorite(isFavorite bool) *CreateTaskCommandBuilder {
	b.isFavorite = isFavorite
	return b
}

func (b *CreateTaskCommandBuilder) WithPriority(priority int32) *CreateTaskCommandBuilder {
	b.priority = domaintask.Priority(priority)
	return b
}

func (b *CreateTaskCommandBuilder) WithIcon(icon int32) *CreateTaskCommandBuilder {
	b.icon = domaintask.Icon(icon)
	return b
}

func (b *CreateTaskCommandBuilder) WithDeadline(deadline time.Time) *CreateTaskCommandBuilder {
	b.deadline = mo.Some(deadline)
	return b
}

func (b *CreateTaskCommandBuilder) WithNotifyAt(notifyAt []time.Time) *CreateTaskCommandBuilder {
	b.notifyAt = notifyAt
	return b
}

func (b *CreateTaskCommandBuilder) Build() *CreateTaskCommand {
	return &CreateTaskCommand{
		OwnerID:     b.ownerID,
		GroupID:     b.groupID,
		Title:       b.title,
		Description: b.description,
		Location:    b.location,
		IsFavorite:  b.isFavorite,
		Priority:    b.priority,
		Icon:        b.icon,
		Deadline:    b.deadline,
		NotifyAt:    b.notifyAt,
	}
}

type CreateHandler interface {
	Handle(ctx context.Context, command *CreateTaskCommand) (uuid.UUID, error)
}

type createHandler struct {
	repository            domaintask.Repository
	accessGuard           AccessGuard
	notificationPublisher domaintask.NotificationEventPublisher
	createdPublisher      domaintask.CreatedTaskPublisher
}

func (h *createHandler) Handle(ctx context.Context, command *CreateTaskCommand) (uuid.UUID, error) {
	if command == nil {
		return uuid.Nil, errors.New("command is missing")
	}

	if groupID, ok := command.GroupID.Get(); ok {
		if err := h.accessGuard.ValidateGroup(ctx, command.OwnerID, groupID, PermissionCreate); err != nil {
			return uuid.Nil, fmt.Errorf("access guard: %w", err)
		}
	}

	b := domaintask.NewBuilder().
		WithOwnerID(command.OwnerID).
		WithTitle(command.Title).
		WithDescription(command.Description).
		WithIsFavorite(command.IsFavorite).
		WithPriority(command.Priority).
		WithIcon(command.Icon).
		WithNotifyAt(command.NotifyAt)

	if value, ok := command.GroupID.Get(); ok {
		b.WithGroupID(value)
	}

	if value, ok := command.Location.Get(); ok {
		b.WithLocation(value)
	}

	if value, ok := command.Deadline.Get(); ok {
		b.WithDeadline(value)
	}

	task, err := b.Build()
	if err != nil {
		return uuid.Nil, fmt.Errorf("build task: %w", err)
	}

	if err = h.repository.Create(ctx, task); err != nil {
		return uuid.Nil, fmt.Errorf("create task: %w", err)
	}

	// TODO: This implementation is intended only for test purposes in MVP stage.
	// TODO: Change this configuration to be atomic: create + notify
	// TODO: Otherwise it failes.
	// FIXME: Change as fast as possible.
	notificationEvent := domaintask.NewNotificationEvent(task.ID, task.OwnerID, task.Title, task.NotifyAt)
	_ = h.notificationPublisher.Notify(ctx, notificationEvent)

	eventBuilder := domaintask.NewCreatedTaskEventBuilder().
		WithTaskID(task.ID).
		WithUserID(task.OwnerID).
		WithCreatedAt(task.CreatedAt)

	if groupID, ok := task.GroupID.Get(); ok {
		eventBuilder.WithGroupID(groupID)
	}

	if location, ok := task.Location.Get(); ok {
		eventBuilder.WithLocation(location)
	}

	createdEvent := eventBuilder.Build()
	_ = h.createdPublisher.Publish(ctx, createdEvent)

	return task.ID, nil
}

func NewCreateHandler(
	repository domaintask.Repository,
	guard AccessGuard,
	notificationPublisher domaintask.NotificationEventPublisher,
	createdPublisher domaintask.CreatedTaskPublisher,
) (CreateHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	if notificationPublisher == nil {
		return nil, errors.New("publisher is missing")
	}

	if createdPublisher == nil {
		return nil, errors.New("created publisher is missing")
	}

	return &createHandler{
		repository:            repository,
		accessGuard:           guard,
		notificationPublisher: notificationPublisher,
		createdPublisher:      createdPublisher,
	}, nil
}

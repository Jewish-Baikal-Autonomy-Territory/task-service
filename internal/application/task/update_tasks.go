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

type UpdateTaskCommand struct {
	ID          uuid.UUID
	RequesterID uuid.UUID
	GroupID     mo.Option[uuid.UUID]

	Title       mo.Option[string]
	Description mo.Option[string]
	Location    mo.Option[domaintask.GeoPoint]
	IsFavorite  mo.Option[bool]
	Priority    mo.Option[domaintask.Priority]
	Icon        mo.Option[domaintask.Icon]
	Deadline    mo.Option[time.Time]
	NotifyAt    mo.Option[[]time.Time]
}

type UpdateTaskCommandBuilder struct {
	id          uuid.UUID
	requesterID uuid.UUID
	groupID     mo.Option[uuid.UUID]

	title       mo.Option[string]
	description mo.Option[string]
	location    mo.Option[domaintask.GeoPoint]
	isFavorite  mo.Option[bool]
	priority    mo.Option[domaintask.Priority]
	icon        mo.Option[domaintask.Icon]
	deadline    mo.Option[time.Time]
	notifyAt    mo.Option[[]time.Time]
}

func (b *UpdateTaskCommandBuilder) WithID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.id = id
	return b
}

func (b *UpdateTaskCommandBuilder) WithRequesterID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.requesterID = id
	return b
}

func (b *UpdateTaskCommandBuilder) WithGroupID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *UpdateTaskCommandBuilder) WithTitle(title string) *UpdateTaskCommandBuilder {
	b.title = mo.Some(title)
	return b
}

func (b *UpdateTaskCommandBuilder) WithDescription(description string) *UpdateTaskCommandBuilder {
	b.description = mo.Some(description)
	return b
}

func (b *UpdateTaskCommandBuilder) WithLocation(location domaintask.GeoPoint) *UpdateTaskCommandBuilder {
	b.location = mo.Some(location)
	return b
}

func (b *UpdateTaskCommandBuilder) WithIsFavorite(isFavorite bool) *UpdateTaskCommandBuilder {
	b.isFavorite = mo.Some(isFavorite)
	return b
}

func (b *UpdateTaskCommandBuilder) WithPriority(priority domaintask.Priority) *UpdateTaskCommandBuilder {
	b.priority = mo.Some(priority)
	return b
}

func (b *UpdateTaskCommandBuilder) WithIcon(icon domaintask.Icon) *UpdateTaskCommandBuilder {
	b.icon = mo.Some(icon)
	return b
}

func (b *UpdateTaskCommandBuilder) WithDeadline(deadline time.Time) *UpdateTaskCommandBuilder {
	b.deadline = mo.Some(deadline)
	return b
}

func (b *UpdateTaskCommandBuilder) WithNotifyAt(notifyAt []time.Time) *UpdateTaskCommandBuilder {
	b.notifyAt = mo.Some(notifyAt)
	return b
}

func (b *UpdateTaskCommandBuilder) Build() (*UpdateTaskCommand, error) {
	if b.id == uuid.Nil {
		return nil, domaintask.ErrInvalidData
	}

	if b.requesterID == uuid.Nil {
		return nil, domaintask.ErrInvalidData
	}

	if value, ok := b.icon.Get(); ok && !value.Valid() {
		return nil, domaintask.ErrInvalidData
	}

	return &UpdateTaskCommand{
		ID:          b.id,
		RequesterID: b.requesterID,
		GroupID:     b.groupID,
		Title:       b.title,
		Description: b.description,
		Location:    b.location,
		Priority:    b.priority,
		Icon:        b.icon,
		IsFavorite:  b.isFavorite,
		Deadline:    b.deadline,
		NotifyAt:    b.notifyAt,
	}, nil
}

func NewUpdateTaskCommandBuilder() *UpdateTaskCommandBuilder {
	return &UpdateTaskCommandBuilder{}
}

type UpdateHandler interface {
	Handle(ctx context.Context, command *UpdateTaskCommand) error
}

type updateHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *updateHandler) Handle(ctx context.Context, command *UpdateTaskCommand) error {
	if command == nil {
		return errors.New("command is missing")
	}

	t, err := h.repository.FindByID(ctx, command.ID)
	if err != nil {
		return err
	}

	if groupID, ok := t.GroupID.Get(); ok {
		if err = h.accessGuard.ValidateGroup(ctx, command.RequesterID, groupID, PermissionUpdate); err != nil {
			return fmt.Errorf("access guard: %w", err)
		}
	} else {
		if err = h.accessGuard.ValidatePersonal(ctx, command.RequesterID, t.OwnerID); err != nil {
			return fmt.Errorf("access guard: %w", err)
		}
	}

	if value, ok := command.Title.Get(); ok {
		t.Title = value
	}

	if value, ok := command.Description.Get(); ok {
		t.Description = value
	}

	if value, ok := command.Location.Get(); ok {
		t.Location = mo.Some(value)
	}

	if value, ok := command.IsFavorite.Get(); ok {
		t.IsFavorite = value
	}

	if value, ok := command.Priority.Get(); ok {
		t.Priority = value
	}

	if value, ok := command.Icon.Get(); ok {
		t.Icon = value
	}

	if value, ok := command.Deadline.Get(); ok {
		t.Deadline = mo.Some(value)
	}

	if value, ok := command.NotifyAt.Get(); ok {
		t.NotifyAt = value
	}

	return h.repository.Save(ctx, t)
}

func NewUpdateHandler(repository domaintask.Repository, guard AccessGuard) (UpdateHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &updateHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

package task

import (
	"errors"
	"task-service/internal/domain/task"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type UpdateTaskCommand struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	GroupID mo.Option[uuid.UUID]

	Title       mo.Option[string]
	Description mo.Option[string]
	Priority    mo.Option[task.Priority]
	IsFavorite  mo.Option[bool]
}

type UpdateTaskCommandBuilder struct {
	id          uuid.UUID
	ownerID     uuid.UUID
	groupID     mo.Option[uuid.UUID]
	title       mo.Option[string]
	description mo.Option[string]
	priority    mo.Option[task.Priority]
	isFavorite  mo.Option[bool]
}

func (b *UpdateTaskCommandBuilder) SetID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.id = id
	return b
}

func (b *UpdateTaskCommandBuilder) SetOwnerID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.ownerID = id
	return b
}

func (b *UpdateTaskCommandBuilder) SetGroupID(id uuid.UUID) *UpdateTaskCommandBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *UpdateTaskCommandBuilder) SetTitle(title string) *UpdateTaskCommandBuilder {
	b.title = mo.Some(title)
	return b
}

func (b *UpdateTaskCommandBuilder) SetDescription(description string) *UpdateTaskCommandBuilder {
	b.description = mo.Some(description)
	return b
}

func (b *UpdateTaskCommandBuilder) SetPriority(priority task.Priority) *UpdateTaskCommandBuilder {
	b.priority = mo.Some(priority)
	return b
}

func (b *UpdateTaskCommandBuilder) SetIsFavorite(isFavorite bool) *UpdateTaskCommandBuilder {
	b.isFavorite = mo.Some(isFavorite)
	return b
}

func (b *UpdateTaskCommandBuilder) Build() (*UpdateTaskCommand, error) {
	if b.id == uuid.Nil {
		return nil, errors.New("id is required")
	}
	if b.ownerID == uuid.Nil {
		return nil, errors.New("owner id is required")
	}
	return &UpdateTaskCommand{
		ID:          b.id,
		OwnerID:     b.ownerID,
		GroupID:     b.groupID,
		Title:       b.title,
		Description: b.description,
		Priority:    b.priority,
		IsFavorite:  b.isFavorite,
	}, nil
}

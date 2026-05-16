package task

import (
	"errors"
	"task-service/internal/domain/task"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type CreateTaskCommand struct {
	OwnerID uuid.UUID
	GroupID mo.Option[uuid.UUID]

	Title       string
	Description string
	IsFavorite  bool
	Priority    task.Priority
}

type CreateTaskCommandBuilder struct {
	ownerID uuid.UUID
	groupID mo.Option[uuid.UUID]

	title       string
	description string
	isFavorite  bool
	priority    task.Priority
}

func NewCreateTaskCommandBuilder() *CreateTaskCommandBuilder {
	return &CreateTaskCommandBuilder{}
}

func (b *CreateTaskCommandBuilder) SetOwnerID(id uuid.UUID) *CreateTaskCommandBuilder {
	b.ownerID = id
	return b
}

func (b *CreateTaskCommandBuilder) SetGroupID(id uuid.UUID) *CreateTaskCommandBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *CreateTaskCommandBuilder) SetTitle(title string) *CreateTaskCommandBuilder {
	b.title = title
	return b
}

func (b *CreateTaskCommandBuilder) SetDescription(description string) *CreateTaskCommandBuilder {
	b.description = description
	return b
}

func (b *CreateTaskCommandBuilder) SetIsFavorite(isFavorite bool) *CreateTaskCommandBuilder {
	b.isFavorite = isFavorite
	return b
}

func (b *CreateTaskCommandBuilder) SetPriority(priority task.Priority) *CreateTaskCommandBuilder {
	b.priority = priority
	return b
}

func (b *CreateTaskCommandBuilder) Build() (*CreateTaskCommand, error) {
	if b.ownerID == uuid.Nil {
		return nil, errors.New("owner id is required")
	}
	if b.title == "" {
		return nil, errors.New("title is required")
	}
	if !b.priority.Valid() {
		return nil, errors.New("priority is invalid")
	}
	return &CreateTaskCommand{
		OwnerID:     b.ownerID,
		GroupID:     b.groupID,
		Title:       b.title,
		Description: b.description,
		IsFavorite:  b.isFavorite,
		Priority:    b.priority,
	}, nil
}

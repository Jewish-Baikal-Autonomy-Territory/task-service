package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type FindDeletedFilter struct {
	OwnerID mo.Option[uuid.UUID]
	GroupID mo.Option[uuid.UUID]
	Cursor  Cursor
}

type FindDeletedFilterBuilder struct {
	ownerID mo.Option[uuid.UUID]
	groupID mo.Option[uuid.UUID]
	cursor  Cursor
}

func (b *FindDeletedFilterBuilder) WithOwnerID(id uuid.UUID) *FindDeletedFilterBuilder {
	b.ownerID = mo.Some(id)
	return b
}

func (b *FindDeletedFilterBuilder) WithGroupID(id uuid.UUID) *FindDeletedFilterBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *FindDeletedFilterBuilder) WithCursor(cursor Cursor) *FindDeletedFilterBuilder {
	b.cursor = cursor
	return b
}

func (b *FindDeletedFilterBuilder) Build() (FindDeletedFilter, error) {
	if b.ownerID.IsNone() && b.groupID.IsNone() {
		return FindDeletedFilter{}, ErrFailedPrecondition
	}

	if b.ownerID.IsSome() && b.ownerID.MustGet() == uuid.Nil {
		return FindDeletedFilter{}, ErrInvalidData
	}

	if b.groupID.IsSome() && b.groupID.MustGet() == uuid.Nil {
		return FindDeletedFilter{}, ErrInvalidData
	}

	return FindDeletedFilter{
		OwnerID: b.ownerID,
		GroupID: b.groupID,
		Cursor:  b.cursor,
	}, nil
}

func NewFindDeletedFilterBuilder() FindDeletedFilterBuilder {
	return FindDeletedFilterBuilder{}
}

type Repository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	FindDeleted(ctx context.Context, filter FindDeletedFilter) ([]uuid.UUID, time.Time, error)
	Save(ctx context.Context, task *Task) error
	SoftDelete(ctx context.Context, task *Task) error
	HardDeleteByUserID(ctx context.Context, id uuid.UUID) error
}

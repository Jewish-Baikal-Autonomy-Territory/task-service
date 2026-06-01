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

type ListDeletedTasksQuery struct {
	RequesterID uuid.UUID
	Cursor      SearchCursor
	GroupID     mo.Option[uuid.UUID]
}

func (q *ListDeletedTasksQuery) Valid() bool {
	if groupID, ok := q.GroupID.Get(); ok && groupID == uuid.Nil {
		return false
	}

	return q.RequesterID != uuid.Nil
}

type ListDeletedTasksQueryBuilder struct {
	requesterID uuid.UUID
	cursor      SearchCursor
	groupID     mo.Option[uuid.UUID]
}

func (b *ListDeletedTasksQueryBuilder) WithRequesterID(id uuid.UUID) *ListDeletedTasksQueryBuilder {
	b.requesterID = id
	return b
}

func (b *ListDeletedTasksQueryBuilder) WithGroupID(id uuid.UUID) *ListDeletedTasksQueryBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *ListDeletedTasksQueryBuilder) WithCursor(cursor SearchCursor) *ListDeletedTasksQueryBuilder {
	b.cursor = cursor
	return b
}

func (b *ListDeletedTasksQueryBuilder) Build() (ListDeletedTasksQuery, error) {
	if b.requesterID == uuid.Nil {
		return ListDeletedTasksQuery{}, domaintask.ErrInvalidData
	}

	if groupID, ok := b.groupID.Get(); ok && groupID == uuid.Nil {
		return ListDeletedTasksQuery{}, domaintask.ErrInvalidData
	}

	return ListDeletedTasksQuery{
		RequesterID: b.requesterID,
		Cursor:      b.cursor,
		GroupID:     b.groupID,
	}, nil
}

func NewListDeletedTasksQueryBuilder() ListDeletedTasksQueryBuilder {
	return ListDeletedTasksQueryBuilder{}
}

type ListDeletedHandler interface {
	Handle(ctx context.Context, query ListDeletedTasksQuery) ([]uuid.UUID, time.Time, error)
}

type listDeletedHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *listDeletedHandler) Handle(ctx context.Context, query ListDeletedTasksQuery) ([]uuid.UUID, time.Time, error) {
	if !query.Valid() {
		return nil, time.Time{}, errors.New("invalid command")
	}

	if groupID, ok := query.GroupID.Get(); ok {
		if err := h.accessGuard.ValidateGroup(ctx, query.RequesterID, groupID, PermissionReadDeleted); err != nil {
			return nil, time.Time{}, fmt.Errorf("access guard: %w", err)
		}
	}

	b := domaintask.NewFindDeletedFilterBuilder()
	b.WithCursor(query.Cursor.toDomain())

	if groupID, ok := query.GroupID.Get(); ok {
		b.WithGroupID(groupID)
	} else {
		b.WithOwnerID(query.RequesterID)
	}

	q, err := b.Build()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("build query: %w", err)
	}

	ids, key, err := h.repository.FindDeleted(ctx, q)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("find deleted: %w", err)
	}

	return ids, key, nil
}

func NewListDeletedHandler(repository domaintask.Repository, guard AccessGuard) (ListDeletedHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &listDeletedHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

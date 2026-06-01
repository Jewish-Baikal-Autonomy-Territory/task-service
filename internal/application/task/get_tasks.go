package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"

	"github.com/google/uuid"
)

type GetQuery struct {
	ID          uuid.UUID
	RequesterID uuid.UUID
}

func (q *GetQuery) Valid() bool {
	return q.ID != uuid.Nil && q.RequesterID != uuid.Nil
}

func NewFindTaskQuery(id, requesterID uuid.UUID) (GetQuery, error) {
	if id == uuid.Nil || requesterID == uuid.Nil {
		return GetQuery{}, domaintask.ErrInvalidData
	}

	return GetQuery{
		ID:          id,
		RequesterID: requesterID,
	}, nil
}

type GetHandler interface {
	Handle(ctx context.Context, query GetQuery) (*domaintask.Task, error)
}

type getHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *getHandler) Handle(ctx context.Context, query GetQuery) (*domaintask.Task, error) {
	if !query.Valid() {
		return nil, errors.New("invalid query")
	}

	task, err := h.repository.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find task: %w", err)
	}

	if groupID, ok := task.GroupID.Get(); ok {
		if err = h.accessGuard.ValidateGroup(ctx, query.RequesterID, groupID, PermissionRead); err != nil {
			return nil, fmt.Errorf("access guard: %w", err)
		}
	} else {
		if err = h.accessGuard.ValidatePersonal(ctx, query.RequesterID, task.OwnerID); err != nil {
			return nil, fmt.Errorf("access guard: %w", err)
		}
	}

	if task.IsDeleted() {
		return nil, domaintask.ErrAlreadyDeleted
	}

	return task, nil
}

func NewGetHandler(repository domaintask.Repository, guard AccessGuard) (GetHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &getHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

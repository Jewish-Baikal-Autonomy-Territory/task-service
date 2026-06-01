package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	"github.com/google/uuid"
)

type DeleteTaskCommand struct {
	ID          uuid.UUID
	RequesterID uuid.UUID
}

func (c *DeleteTaskCommand) Valid() bool {
	return c.ID != uuid.Nil && c.RequesterID != uuid.Nil
}

func NewDeleteTaskCommand(id, requesterID uuid.UUID) (DeleteTaskCommand, error) {
	if id == uuid.Nil || requesterID == uuid.Nil {
		return DeleteTaskCommand{}, domaintask.ErrInvalidData
	}

	return DeleteTaskCommand{
		ID:          id,
		RequesterID: requesterID,
	}, nil
}

type DeleteResult struct {
	ID      uuid.UUID
	PurgeAt time.Time
}

func newDeleteResult(id uuid.UUID, purgeAt time.Time) (DeleteResult, error) {
	if id == uuid.Nil {
		return DeleteResult{}, domaintask.ErrInvalidData
	}

	return DeleteResult{
		ID:      id,
		PurgeAt: purgeAt,
	}, nil
}

type DeleteHandler interface {
	Handle(ctx context.Context, command DeleteTaskCommand) (DeleteResult, error)
}

type deleteHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *deleteHandler) Handle(ctx context.Context, command DeleteTaskCommand) (DeleteResult, error) {
	if !command.Valid() {
		return DeleteResult{}, errors.New("invalid command")
	}

	task, err := h.repository.FindByID(ctx, command.ID)
	if err != nil {
		return DeleteResult{}, fmt.Errorf("find task: %w", err)
	}

	if groupID, ok := task.GroupID.Get(); ok {
		if err = h.accessGuard.ValidateGroup(ctx, command.RequesterID, groupID, PermissionDelete); err != nil {
			return DeleteResult{}, fmt.Errorf("access guard: %w", err)
		}
	} else {
		if err = h.accessGuard.ValidatePersonal(ctx, command.RequesterID, task.OwnerID); err != nil {
			return DeleteResult{}, fmt.Errorf("access guard: %w", err)
		}
	}

	if err = task.SoftDelete(); err != nil {
		return DeleteResult{}, fmt.Errorf("soft delete task: %w", err)
	}

	if err = h.repository.Save(ctx, task); err != nil {
		return DeleteResult{}, fmt.Errorf("update task: %w", err)
	}

	res, err := newDeleteResult(task.ID, task.PurgeAt.MustGet())
	if err != nil {
		return DeleteResult{}, fmt.Errorf("create delete result: %w", err)
	}

	return res, nil
}

func NewDeleteHandler(repository domaintask.Repository, guard AccessGuard) (DeleteHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &deleteHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

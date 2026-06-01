package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"

	"github.com/google/uuid"
)

type RestoreTaskCommand struct {
	ID          uuid.UUID
	RequesterID uuid.UUID
}

func (c *RestoreTaskCommand) Valid() bool {
	return c.ID != uuid.Nil && c.RequesterID != uuid.Nil
}

func NewRestoreTaskCommand(id, requesterID uuid.UUID) (RestoreTaskCommand, error) {
	if id == uuid.Nil || requesterID == uuid.Nil {
		return RestoreTaskCommand{}, domaintask.ErrInvalidData
	}

	return RestoreTaskCommand{
		ID:          id,
		RequesterID: requesterID,
	}, nil
}

type RestoreHandler interface {
	Handle(ctx context.Context, command RestoreTaskCommand) error
}

type restoreHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *restoreHandler) Handle(ctx context.Context, command RestoreTaskCommand) error {
	if !command.Valid() {
		return errors.New("invalid command")
	}

	task, err := h.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}

	if groupID, ok := task.GroupID.Get(); ok {
		if err = h.accessGuard.ValidateGroup(ctx, command.RequesterID, groupID, PermissionRestore); err != nil {
			return fmt.Errorf("access guard: %w", err)
		}
	} else {
		if err = h.accessGuard.ValidatePersonal(ctx, command.RequesterID, task.OwnerID); err != nil {
			return fmt.Errorf("access guard: %w", err)
		}
	}

	if err = task.Restore(); err != nil {
		return fmt.Errorf("restore task: %w", err)
	}

	if err = h.repository.Save(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func NewRestoreHandler(repository domaintask.Repository, guard AccessGuard) (RestoreHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &restoreHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

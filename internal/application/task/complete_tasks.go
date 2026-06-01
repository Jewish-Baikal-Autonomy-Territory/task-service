package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"

	"github.com/google/uuid"
)

type CompleteTaskCommand struct {
	ID          uuid.UUID
	RequesterID uuid.UUID
}

func (c *CompleteTaskCommand) Valid() bool {
	return c.ID != uuid.Nil && c.RequesterID != uuid.Nil
}

func NewCompleteTaskCommand(id, requesterID uuid.UUID) (CompleteTaskCommand, error) {
	if id == uuid.Nil || requesterID == uuid.Nil {
		return CompleteTaskCommand{}, domaintask.ErrInvalidData
	}

	return CompleteTaskCommand{
		ID:          id,
		RequesterID: requesterID,
	}, nil
}

type CompleteHandler interface {
	Handle(ctx context.Context, command CompleteTaskCommand) error
}

type completeHandler struct {
	repository  domaintask.Repository
	accessGuard AccessGuard
}

func (h *completeHandler) Handle(ctx context.Context, command CompleteTaskCommand) error {
	if !command.Valid() {
		return errors.New("invalid command")
	}

	task, err := h.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}

	if groupID, ok := task.GroupID.Get(); ok {
		if err = h.accessGuard.ValidateGroup(ctx, command.RequesterID, groupID, PermissionComplete); err != nil {
			return fmt.Errorf("validate group: %w", err)
		}
	} else {
		if err = h.accessGuard.ValidatePersonal(ctx, command.RequesterID, task.OwnerID); err != nil {
			return fmt.Errorf("validate personal: %w", err)
		}
	}

	if err = task.Complete(); err != nil {
		return fmt.Errorf("complete task: %w", err)
	}

	if err = h.repository.Save(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func NewCompleteHandler(repository domaintask.Repository, guard AccessGuard) (CompleteHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &completeHandler{
		repository:  repository,
		accessGuard: guard,
	}, nil
}

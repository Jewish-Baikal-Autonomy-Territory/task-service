package event

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"

	"github.com/google/uuid"
)

type DeleteUserTasksHandler interface {
	Handle(ctx context.Context, userID uuid.UUID) error
}
type deleteUserTasksHandler struct {
	repository domaintask.Repository
}

func (h *deleteUserTasksHandler) Handle(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return domaintask.ErrInvalidData
	}

	if err := h.repository.HardDeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("delete user tasks: %w", err)
	}

	return nil
}

func NewDeleteUserTasksHandler(repository domaintask.Repository) (DeleteUserTasksHandler, error) {
	if repository == nil {
		return nil, errors.New("repository is missing")
	}

	return &deleteUserTasksHandler{
		repository: repository,
	}, nil
}

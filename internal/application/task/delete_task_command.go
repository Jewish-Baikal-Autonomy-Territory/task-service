package task

import (
	"errors"

	"github.com/google/uuid"
)

type DeleteTaskCommand struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}

func NewDeleteTaskCommand(id, ownerID uuid.UUID) (DeleteTaskCommand, error) {
	if id == uuid.Nil {
		return DeleteTaskCommand{}, errors.New("id is required")
	}
	if ownerID == uuid.Nil {
		return DeleteTaskCommand{}, errors.New("owner id is required")
	}
	return DeleteTaskCommand{
		ID:      id,
		OwnerID: ownerID,
	}, nil
}

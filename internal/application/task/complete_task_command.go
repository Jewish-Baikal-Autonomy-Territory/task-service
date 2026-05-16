package task

import (
	"errors"

	"github.com/google/uuid"
)

type CompleteTaskCommand struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}

func NewCompleteTaskCommand(id, ownerID uuid.UUID) (CompleteTaskCommand, error) {
	if id == uuid.Nil {
		return CompleteTaskCommand{}, errors.New("id is required")
	}
	if ownerID == uuid.Nil {
		return CompleteTaskCommand{}, errors.New("owner id is required")
	}
	return CompleteTaskCommand{
		ID:      id,
		OwnerID: ownerID,
	}, nil
}

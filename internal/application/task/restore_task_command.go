package task

import (
	"errors"

	"github.com/google/uuid"
)

type RestoreTaskCommand struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}

func NewRestoreTaskCommand(id, ownerID uuid.UUID) (RestoreTaskCommand, error) {
	if id == uuid.Nil {
		return RestoreTaskCommand{}, errors.New("id is required")
	}
	if ownerID == uuid.Nil {
		return RestoreTaskCommand{}, errors.New("owner is required")
	}
	return RestoreTaskCommand{
		ID:      id,
		OwnerID: ownerID,
	}, nil
}

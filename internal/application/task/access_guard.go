package task

import (
	"context"
	"errors"
	domaintask "task-service/internal/domain/task"

	"github.com/google/uuid"
)

var ErrPermissionDenied = errors.New("permission denied")

type Permission int32

const (
	PermissionUnknown Permission = iota
	PermissionCreate             = 1 << iota
	PermissionUpdate
	PermissionDelete
	PermissionRestore
	PermissionRead
	PermissionReadDeleted
	PermissionComplete
)

var permStrMapper = map[Permission]string{
	PermissionCreate:      "task:create",
	PermissionUpdate:      "task:update",
	PermissionDelete:      "task:delete",
	PermissionRestore:     "task:restore",
	PermissionRead:        "task:read",
	PermissionReadDeleted: "task:read-deleted",
	PermissionComplete:    "task:complete",
}

func (p Permission) String() string {
	if value, ok := permStrMapper[p]; ok {
		return value
	}
	return "task:unknown"
}

var strPermMapper = map[string]Permission{
	"task:create":       PermissionCreate,
	"task:update":       PermissionUpdate,
	"task:delete":       PermissionDelete,
	"task:restore":      PermissionDelete,
	"task:read":         PermissionRead,
	"task:read-deleted": PermissionReadDeleted,
	"task:complete":     PermissionComplete,
}

func ParsePermission(s string) (Permission, error) {
	if value, ok := strPermMapper[s]; ok {
		return value, nil
	}
	return PermissionUnknown, domaintask.ErrInvalidData
}

type AccessGuard interface {
	ValidatePersonal(ctx context.Context, requesterID, taskOwnerID uuid.UUID) error
	ValidateGroup(ctx context.Context, requesterID, groupID uuid.UUID, perms ...Permission) error
}

type accessGuard struct {
}

func (g *accessGuard) ValidatePersonal(ctx context.Context, requesterID, taskOwnerID uuid.UUID) error {
	if requesterID != taskOwnerID {
		return ErrPermissionDenied
	}
	return nil
}

func (g *accessGuard) ValidateGroup(ctx context.Context, requesterID, groupID uuid.UUID, perms ...Permission) error {
	return ErrUnsupported
}

func NewAccessGuard() AccessGuard {
	return &accessGuard{}
}

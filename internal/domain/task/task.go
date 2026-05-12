package task

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

var (
	ErrAlreadyDeleted         = errors.New("task deleted")
	ErrNotDeleted             = errors.New("task is not deleted")
	ErrRestoreWindowClosed    = errors.New("task is not deleted")
	ErrInvalidTaskTitle       = errors.New("invalid task title")
	ErrInvalidTaskDescription = errors.New("invalid task description")
	ErrInvalidTaskPriority    = errors.New("invalid task priority")
)

const restoreWindow = 14 * 24 * time.Hour

type Task struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	GroupID mo.Option[uuid.UUID]

	Title       string
	Description string
	Location    mo.Option[GeoPoint]
	IsFavorite  bool
	Priority    Priority
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt mo.Option[time.Time]
	Deadline    mo.Option[time.Time]
	PurgeAt     mo.Option[time.Time]
}

func (t *Task) IsDeleted() bool {
	return t.PurgeAt.IsSome()
}

func (t *Task) IsCompleted() bool {
	return t.CompletedAt.IsSome()
}

func (t *Task) Complete() error {
	if t.IsDeleted() {
		return ErrAlreadyDeleted
	}

	currentTime := time.Now()
	t.CompletedAt = mo.Some(currentTime)
	t.UpdatedAt = currentTime
	return nil
}

func (t *Task) SoftDelete() error {
	if t.IsDeleted() {
		return ErrAlreadyDeleted
	}

	currentTime := time.Now()
	t.UpdatedAt = currentTime
	t.PurgeAt = mo.Some(currentTime.Add(restoreWindow))
	return nil
}

func (t *Task) CanRestore() bool {
	purgeAt, ok := t.PurgeAt.Get()
	if !ok {
		return false
	}
	return time.Now().Before(purgeAt)
}

func (t *Task) Restore() error {
	if !t.IsDeleted() {
		return ErrNotDeleted
	}
	if !t.CanRestore() {
		return ErrRestoreWindowClosed
	}

	t.UpdatedAt = time.Now()
	t.PurgeAt = mo.None[time.Time]()
	return nil
}

func NewTask(ownerID uuid.UUID, title, description string, priority Priority) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTaskTitle
	}
	if strings.TrimSpace(description) == "" {
		return nil, ErrInvalidTaskDescription
	}
	if !priority.Valid() {
		return nil, ErrInvalidTaskPriority
	}

	currentTime := time.Now()
	return &Task{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      StatusPending,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}, nil
}

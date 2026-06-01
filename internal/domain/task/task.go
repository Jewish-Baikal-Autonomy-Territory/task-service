package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

var (
	ErrAlreadyDeleted      = errors.New("task deleted")
	ErrNotDeleted          = errors.New("task is not deleted")
	ErrRestoreWindowClosed = errors.New("task is not deleted")
	ErrInvalidData         = errors.New("invalid data")
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
	Icon        Icon
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt mo.Option[time.Time]
	Deadline    mo.Option[time.Time]
	PurgeAt     mo.Option[time.Time]
	NotifyAt    []time.Time
}

func (t *Task) IsDeleted() bool {
	return t.PurgeAt.IsSome()
}

func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}

func (t *Task) Complete() error {
	if t.IsDeleted() {
		return ErrAlreadyDeleted
	}

	t.Status = StatusCompleted
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

type Builder struct {
	ownerID uuid.UUID
	groupID mo.Option[uuid.UUID]

	title       string
	description string
	location    mo.Option[GeoPoint]
	isFavorite  bool
	priority    Priority
	icon        Icon
	deadline    mo.Option[time.Time]
	notifyAt    []time.Time
}

func (b *Builder) WithOwnerID(ownerID uuid.UUID) *Builder {
	b.ownerID = ownerID
	return b
}

func (b *Builder) WithGroupID(groupID uuid.UUID) *Builder {
	b.groupID = mo.Some(groupID)
	return b
}

func (b *Builder) WithTitle(title string) *Builder {
	b.title = title
	return b
}

func (b *Builder) WithDescription(description string) *Builder {
	b.description = description
	return b
}

func (b *Builder) WithLocation(location GeoPoint) *Builder {
	b.location = mo.Some(location)
	return b
}

func (b *Builder) WithIsFavorite(isFavorite bool) *Builder {
	b.isFavorite = isFavorite
	return b
}

func (b *Builder) WithPriority(priority Priority) *Builder {
	b.priority = priority
	return b
}

func (b *Builder) WithIcon(icon Icon) *Builder {
	b.icon = icon
	return b
}

func (b *Builder) WithDeadline(deadline time.Time) *Builder {
	b.deadline = mo.Some(deadline)
	return b
}

func (b *Builder) WithNotifyAt(notifyAt []time.Time) *Builder {
	b.notifyAt = notifyAt
	return b
}

func (b *Builder) Build() (*Task, error) {
	if b.ownerID == uuid.Nil {
		return nil, ErrInvalidData
	}

	if b.groupID.IsSome() && b.groupID.MustGet() == uuid.Nil {
		return nil, ErrInvalidData
	}

	if b.title == "" {
		return nil, ErrInvalidData
	}

	if b.description == "" {
		return nil, ErrInvalidData
	}

	if !b.priority.Valid() {
		return nil, ErrInvalidData
	}

	if !b.icon.Valid() {
		return nil, ErrInvalidData
	}

	currentTime := time.Now()

	deadline, ok := b.deadline.Get()
	if ok && deadline.Before(currentTime) {
		return nil, ErrInvalidData
	}

	for _, notifyAt := range b.notifyAt {
		if notifyAt.Before(currentTime) {
			return nil, ErrInvalidData
		}
		if ok && deadline.Before(notifyAt) {
			return nil, ErrInvalidData
		}
	}

	return &Task{
		ID:          uuid.New(),
		OwnerID:     b.ownerID,
		GroupID:     b.groupID,
		Title:       b.title,
		Description: b.description,
		Location:    b.location,
		IsFavorite:  b.isFavorite,
		Priority:    b.priority,
		Icon:        b.icon,
		Status:      StatusPending,
		Deadline:    b.deadline,
		NotifyAt:    b.notifyAt,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}, nil
}

func NewBuilder() *Builder {
	return &Builder{}
}

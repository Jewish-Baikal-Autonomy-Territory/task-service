package mongo

import (
	"task-service/internal/domain/task"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type taskDocument struct {
	ID      bson.Binary  `bson:"_id"`
	OwnerID bson.Binary  `bson:"owner_id"`
	GroupID *bson.Binary `bson:"group_id,omitempty"`

	Title       string            `bson:"title"`
	Description string            `bson:"description"`
	Location    *geoPointDocument `bson:"location,omitempty"`
	IsFavorite  bool              `bson:"is_favorite"`
	Status      string            `bson:"status"`
	Priority    string            `bson:"priority"`
	CreatedAt   time.Time         `bson:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"`
	CompletedAt *time.Time        `bson:"completed_at,omitempty"`
	Deadline    *time.Time        `bson:"deadline,omitempty"`
	PurgeAt     *time.Time        `bson:"purge_at,omitempty"`
}

func (t *taskDocument) toDomain() *task.Task {
	id, _ := uuid.FromBytes(t.ID.Data)
	ownerID, _ := uuid.FromBytes(t.OwnerID.Data)
	domainTask := &task.Task{
		ID:          id,
		OwnerID:     ownerID,
		GroupID:     mo.None[uuid.UUID](),
		Title:       t.Title,
		Description: t.Description,
		Location:    mo.None[task.GeoPoint](),
		IsFavorite:  t.IsFavorite,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: mo.None[time.Time](),
		Deadline:    mo.None[time.Time](),
		PurgeAt:     mo.None[time.Time](),
	}
	if t.GroupID != nil {
		groupID, _ := uuid.FromBytes(t.GroupID.Data)
		domainTask.GroupID = mo.Some(groupID)
	}
	if t.Location != nil {
		domainTask.Location = mo.Some[task.GeoPoint](t.Location.toDomain())
	}
	if priority, err := task.ParsePriority(t.Priority); err == nil {
		domainTask.Priority = priority
	}
	if status, err := task.ParseStatus(t.Status); err == nil {
		domainTask.Status = status
	}
	if t.CompletedAt != nil {
		domainTask.CompletedAt = mo.Some[time.Time](*t.CompletedAt)
	}
	if t.Deadline != nil {
		domainTask.Deadline = mo.Some[time.Time](*t.Deadline)
	}
	if t.PurgeAt != nil {
		domainTask.PurgeAt = mo.Some[time.Time](*t.PurgeAt)
	}
	return domainTask
}

func fromDomainTask(t *task.Task) *taskDocument {
	doc := &taskDocument{
		ID:          bson.Binary{Subtype: 0x04, Data: t.ID[:]},
		OwnerID:     bson.Binary{Subtype: 0x04, Data: t.OwnerID[:]},
		Title:       t.Title,
		Description: t.Description,
		IsFavorite:  t.IsFavorite,
		Status:      t.Status.String(),
		Priority:    t.Priority.String(),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	if value, ok := t.GroupID.Get(); ok {
		doc.GroupID = &bson.Binary{Subtype: 0x04, Data: value[:]}
	}
	if value, ok := t.Location.Get(); ok {
		doc.Location = fromDomainGeoPoint(value)
	}
	if value, ok := t.CompletedAt.Get(); ok {
		doc.CompletedAt = &value
	}
	if value, ok := t.Deadline.Get(); ok {
		doc.Deadline = &value
	}
	if value, ok := t.PurgeAt.Get(); ok {
		doc.PurgeAt = &value
	}
	return doc
}

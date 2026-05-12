package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

var (
	ErrEmptyFilter = errors.New("empty filter")
)

type GeoFilter struct {
	Point  GeoPoint
	Radius float32
}

type Filter struct {
	OwnerID    mo.Option[uuid.UUID]
	GroupID    mo.Option[uuid.UUID]
	IsFavorite mo.Option[bool]
	Status     mo.Option[Status]
	Priority   mo.Option[Priority]
	Area       mo.Option[GeoFilter]
}

func (f *Filter) Validate() error {
	if f.OwnerID.IsNone() && f.GroupID.IsNone() {
		return ErrEmptyFilter
	}
	return nil
}

type FilterBuilder struct {
	ownerID    mo.Option[uuid.UUID]
	groupID    mo.Option[uuid.UUID]
	isFavorite mo.Option[bool]
	status     mo.Option[Status]
	priority   mo.Option[Priority]
	area       mo.Option[GeoFilter]
}

func (fb *FilterBuilder) SetOwnerID(id uuid.UUID) *FilterBuilder {
	fb.ownerID = mo.Some(id)
	return fb
}

func (fb *FilterBuilder) SetGroupID(id uuid.UUID) *FilterBuilder {
	fb.groupID = mo.Some(id)
	return fb
}

func (fb *FilterBuilder) SetIsFavorite(isFavorite bool) *FilterBuilder {
	fb.isFavorite = mo.Some(isFavorite)
	return fb
}

func (fb *FilterBuilder) SetStatus(status Status) *FilterBuilder {
	fb.status = mo.Some(status)
	return fb
}

func (fb *FilterBuilder) SetPriority(priority Priority) *FilterBuilder {
	fb.priority = mo.Some(priority)
	return fb
}

func (fb *FilterBuilder) SetArea(point GeoPoint, radius float32) *FilterBuilder {
	fb.area = mo.Some(GeoFilter{
		Point:  point,
		Radius: radius,
	})
	return fb
}

func (fb *FilterBuilder) Build() (*Filter, error) {
	filter := &Filter{
		OwnerID:    fb.ownerID,
		GroupID:    fb.groupID,
		IsFavorite: fb.isFavorite,
		Status:     fb.status,
		Priority:   fb.priority,
		Area:       fb.area,
	}
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validate filter: %w", err)
	}
	return filter, nil
}

//go:generate mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock
type Repository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	FindAll(ctx context.Context, filter *Filter) ([]*Task, error)
	FindDeleted(ctx context.Context, filter *Filter) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Close(ctx context.Context) error
}

package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type Cursor struct {
	Key   time.Time
	Limit uint64
}

func NewCursor(limit uint64, key time.Time) Cursor {
	return Cursor{
		Key:   key,
		Limit: limit,
	}
}

type Filter struct {
	SearchQuery string
	OwnerID     mo.Option[uuid.UUID]
	GroupID     mo.Option[uuid.UUID]
	Area        mo.Option[GeoFilter]
	IsFavorite  mo.Option[bool]
	Status      mo.Option[Status]
	Priority    mo.Option[Priority]
	Cursor      Cursor
}

func (f *Filter) Validate() error {
	if f.SearchQuery == "" {
		return ErrInvalidData
	}

	if f.OwnerID.IsNone() && f.GroupID.IsNone() {
		return ErrFailedPrecondition
	}

	return nil
}

type FilterBuilder struct {
	searchQuery string
	ownerID     mo.Option[uuid.UUID]
	groupID     mo.Option[uuid.UUID]
	area        mo.Option[GeoFilter]
	isFavorite  mo.Option[bool]
	status      mo.Option[Status]
	priority    mo.Option[Priority]
	cursor      Cursor
}

func (fb *FilterBuilder) WithSearchQuery(searchQuery string) *FilterBuilder {
	fb.searchQuery = searchQuery
	return fb
}

func (fb *FilterBuilder) WithOwnerID(id uuid.UUID) *FilterBuilder {
	fb.ownerID = mo.Some(id)
	return fb
}

func (fb *FilterBuilder) WithGroupID(id uuid.UUID) *FilterBuilder {
	fb.groupID = mo.Some(id)
	return fb
}

func (fb *FilterBuilder) WithArea(point GeoPoint, radius float64) *FilterBuilder {
	fb.area = mo.Some(GeoFilter{
		Point:  point,
		Radius: radius,
	})
	return fb
}

func (fb *FilterBuilder) WithIsFavorite(isFavorite bool) *FilterBuilder {
	fb.isFavorite = mo.Some(isFavorite)
	return fb
}

func (fb *FilterBuilder) WithStatus(status Status) *FilterBuilder {
	fb.status = mo.Some(status)
	return fb
}

func (fb *FilterBuilder) WithPriority(priority Priority) *FilterBuilder {
	fb.priority = mo.Some(priority)
	return fb
}

func (fb *FilterBuilder) WithCursor(cursor Cursor) *FilterBuilder {
	fb.cursor = cursor
	return fb
}

func (fb *FilterBuilder) Build() (*Filter, error) {
	filter := &Filter{
		SearchQuery: fb.searchQuery,
		OwnerID:     fb.ownerID,
		GroupID:     fb.groupID,
		Area:        fb.area,
		IsFavorite:  fb.isFavorite,
		Status:      fb.status,
		Priority:    fb.priority,
		Cursor:      fb.cursor,
	}
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validate filter: %w", err)
	}
	return filter, nil
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{}
}

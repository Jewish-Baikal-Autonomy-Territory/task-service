package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type GeoFilterQuery struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

func (q *GeoFilterQuery) toDomain() (domaintask.GeoFilter, error) {
	point, err := domaintask.NewGeoPoint(q.Latitude, q.Longitude)
	if err != nil {
		return domaintask.GeoFilter{}, fmt.Errorf("geo conv: %w", err)
	}

	return domaintask.NewGeoFilter(point, q.Radius)
}

func NewGeoFilterQuery(latitude, longitude, radius float64) (GeoFilterQuery, error) {
	if latitude < -90.0 || latitude > 90.0 {
		return GeoFilterQuery{}, domaintask.ErrInvalidData
	}

	if longitude < -180.0 || longitude > 180.0 {
		return GeoFilterQuery{}, domaintask.ErrInvalidData
	}

	if radius < 0.0 {
		return GeoFilterQuery{}, domaintask.ErrInvalidData
	}

	return GeoFilterQuery{
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}, nil
}

type SearchCursor struct {
	Key   time.Time
	Limit uint64
}

func (s *SearchCursor) toDomain() domaintask.Cursor {
	return domaintask.Cursor{
		Key:   s.Key,
		Limit: s.Limit,
	}
}

func NewSearchCursor(limit uint64, key time.Time) SearchCursor {
	return SearchCursor{
		Key:   key,
		Limit: limit,
	}
}

type SearchTasksQuery struct {
	SearchQuery string
	RequesterID uuid.UUID
	GroupID     mo.Option[uuid.UUID]
	GeoQuery    mo.Option[GeoFilterQuery]
	IsFavorite  mo.Option[bool]
	Priority    mo.Option[int32]
	Status      mo.Option[int32]
	Cursor      SearchCursor
}

type SearchTasksQueryBuilder struct {
	searchQuery string
	requesterID uuid.UUID
	groupID     mo.Option[uuid.UUID]
	geoQuery    mo.Option[GeoFilterQuery]
	isFavorite  mo.Option[bool]
	priority    mo.Option[int32]
	status      mo.Option[int32]
	cursor      SearchCursor
}

func (b *SearchTasksQueryBuilder) WithSearchQuery(searchQuery string) *SearchTasksQueryBuilder {
	b.searchQuery = searchQuery
	return b
}

func (b *SearchTasksQueryBuilder) WithRequesterID(requesterID uuid.UUID) *SearchTasksQueryBuilder {
	b.requesterID = requesterID
	return b
}

func (b *SearchTasksQueryBuilder) WithGroupID(id uuid.UUID) *SearchTasksQueryBuilder {
	b.groupID = mo.Some(id)
	return b
}

func (b *SearchTasksQueryBuilder) WithGeoQuery(query GeoFilterQuery) *SearchTasksQueryBuilder {
	b.geoQuery = mo.Some(query)
	return b
}

func (b *SearchTasksQueryBuilder) WithIsFavorite(value bool) *SearchTasksQueryBuilder {
	b.isFavorite = mo.Some(value)
	return b
}

func (b *SearchTasksQueryBuilder) WithPriority(value int32) *SearchTasksQueryBuilder {
	b.priority = mo.Some(value)
	return b
}

func (b *SearchTasksQueryBuilder) WithStatus(value int32) *SearchTasksQueryBuilder {
	b.status = mo.Some(value)
	return b
}

func (b *SearchTasksQueryBuilder) WithCursor(cursor SearchCursor) *SearchTasksQueryBuilder {
	b.cursor = cursor
	return b
}

func (b *SearchTasksQueryBuilder) Build() (*SearchTasksQuery, error) {
	if b.searchQuery == "" {
		return nil, domaintask.ErrInvalidData
	}

	if b.requesterID == uuid.Nil {
		return nil, domaintask.ErrInvalidData
	}

	if groupID, ok := b.groupID.Get(); ok && groupID == uuid.Nil {
		return nil, domaintask.ErrInvalidData
	}

	return &SearchTasksQuery{
		SearchQuery: b.searchQuery,
		RequesterID: b.requesterID,
		GroupID:     b.groupID,
		GeoQuery:    b.geoQuery,
		IsFavorite:  b.isFavorite,
		Priority:    b.priority,
		Status:      b.status,
		Cursor:      b.cursor,
	}, nil
}

func NewSearchTasksQueryBuilder() *SearchTasksQueryBuilder {
	return &SearchTasksQueryBuilder{}
}

type SearchHandler interface {
	Handle(ctx context.Context, query *SearchTasksQuery) ([]uuid.UUID, time.Time, error)
}

type searchHandler struct {
	searcher    domaintask.Searcher
	accessGuard AccessGuard
}

func (h *searchHandler) Handle(ctx context.Context, query *SearchTasksQuery) ([]uuid.UUID, time.Time, error) {
	if query == nil {
		return nil, time.Time{}, errors.New("query is missing")
	}

	b := domaintask.NewFilterBuilder().
		WithSearchQuery(query.SearchQuery)

	if groupID, ok := query.GroupID.Get(); ok {
		if err := h.accessGuard.ValidateGroup(ctx, query.RequesterID, groupID, PermissionRead); err != nil {
			return nil, time.Time{}, fmt.Errorf("access guard: %w", err)
		}

		b.WithGroupID(groupID)
	} else {
		b.WithOwnerID(query.RequesterID)
	}

	if value, ok := query.GeoQuery.Get(); ok {
		geoFilter, err := value.toDomain()
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("geo filter: %w", err)
		}

		b.WithArea(geoFilter.Point, geoFilter.Radius)
	}

	if value, ok := query.IsFavorite.Get(); ok {
		b.WithIsFavorite(value)
	}

	if value, ok := query.Priority.Get(); ok {
		b.WithPriority(domaintask.Priority(value))
	}

	if value, ok := query.Status.Get(); ok {
		b.WithStatus(domaintask.Status(value))
	}

	b.WithCursor(query.Cursor.toDomain())

	searchQuery, err := b.Build()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("build search query: %w", err)
	}

	ids, key, err := h.searcher.Search(ctx, searchQuery)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("searcher: %w", err)
	}

	return ids, key, nil
}

func NewSearchHandler(searcher domaintask.Searcher, guard AccessGuard) (SearchHandler, error) {
	if searcher == nil {
		return nil, errors.New("searcher is missing")
	}

	if guard == nil {
		return nil, errors.New("guard is missing")
	}

	return &searchHandler{
		searcher:    searcher,
		accessGuard: guard,
	}, nil
}

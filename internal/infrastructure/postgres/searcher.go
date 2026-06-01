package postgres

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	sqld "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskSearcher struct {
	client           *pgxpool.Pool
	languageDetector LanguageDetector
}

func (s *TaskSearcher) fromDomainFilter(filter *domaintask.Filter) (string, []any, error) {
	if filter == nil {
		return "", nil, errors.New("no filter found")
	}

	detectedLang := s.languageDetector.Detect(filter.SearchQuery)

	query := sqld.Select("id, created_at").
		From("pgtasks.task").
		Where("created_at < ?", filter.Cursor.Key).
		Where("search_vector @@ websearch_to_tsquery(?, ?)", detectedLang, filter.SearchQuery).
		PlaceholderFormat(sqld.Dollar)

	if value, ok := filter.OwnerID.Get(); ok {
		query = query.Where("owner_id = ?", value.String())
	}

	if value, ok := filter.GroupID.Get(); ok {
		query = query.Where("group_id = ?", value.String())
	}

	if value, ok := filter.Area.Get(); ok {
		query = query.Where(
			"ST_DWithin(location, ST_MakePoint(?, ?), ?)",
			value.Point.Longitude,
			value.Point.Latitude,
			value.Radius,
		)
	}

	if value, ok := filter.IsFavorite.Get(); ok {
		query = query.Where("is_favorite = ?", value)
	}

	if value, ok := filter.Status.Get(); ok {
		query = query.Where("status = ?", value.String())
	}

	if value, ok := filter.Priority.Get(); ok {
		query = query.Where("priority = ?", value.String())
	}

	query = query.Where("purge_at IS NULL").
		OrderByClause(
			"ts_rank_cd(search_vector, websearch_to_tsquery(?, ?)), created_at DESC",
			detectedLang,
			filter.SearchQuery,
		).
		Limit(filter.Cursor.Limit)

	return query.ToSql()
}

func (s *TaskSearcher) Search(ctx context.Context, filter *domaintask.Filter) ([]uuid.UUID, time.Time, error) {
	query, args, err := s.fromDomainFilter(filter)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := s.client.Query(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, time.Time{}, domaintask.ErrNotFound
	}

	if err != nil {
		return nil, time.Time{}, fmt.Errorf("query: %w", err)
	}

	var ids []uuid.UUID
	var createdAt time.Time
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id, &createdAt); err != nil {
			return nil, time.Time{}, fmt.Errorf("scan: %w", err)
		}

		ids = append(ids, id)
	}

	if rows.Err() != nil {
		return nil, time.Time{}, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(ids) == 0 {
		return nil, time.Time{}, domaintask.ErrNotFound
	}

	return ids, createdAt, nil
}

func NewTaskSearcher(client *pgxpool.Pool, detector LanguageDetector) *TaskSearcher {
	return &TaskSearcher{
		client:           client,
		languageDetector: detector,
	}
}

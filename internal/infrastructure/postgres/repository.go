package postgres

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"task-service/internal/infrastructure/postgres/sqlcgen"
	"time"

	sqld "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/mo"
	"github.com/twpayne/go-geom"
)

type TaskRepository struct {
	client   *pgxpool.Pool
	queries  *sqlcgen.Queries
	detector LanguageDetector
}

func fromDomainTime(t time.Time) *time.Time {
	if !t.IsZero() {
		return &t
	}
	return nil
}

func toDomainTask(task *sqlcgen.GetTaskByIDRow) (*domaintask.Task, error) {
	priority, err := domaintask.ParsePriority(task.Priority)
	if err != nil {
		return nil, fmt.Errorf("parse priority: %w", err)
	}

	status, err := domaintask.ParseStatus(task.Status)
	if err != nil {
		return nil, fmt.Errorf("parse status: %w", err)
	}

	icon, err := domaintask.ParseIcon(task.Icon)
	if err != nil {
		return nil, fmt.Errorf("parse icon: %w", err)
	}

	domainTask := &domaintask.Task{
		ID:          task.ID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		IsFavorite:  task.IsFavorite,
		Priority:    priority,
		Icon:        icon,
		Status:      status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.GroupID != nil {
		domainTask.GroupID = mo.Some(*task.GroupID)
	}

	if task.Location != nil {
		lat, long := fromPoint(task.Location)
		geoPoint, err := domaintask.NewGeoPoint(lat, long)
		if err != nil {
			return nil, fmt.Errorf("create geo point: %w", err)
		}
		domainTask.Location = mo.Some(geoPoint)
	}

	if task.CompletedAt != nil {
		domainTask.CompletedAt = mo.Some(*task.CompletedAt)
	}

	if task.Deadline != nil {
		domainTask.Deadline = mo.Some(*task.Deadline)
	}

	if task.PurgeAt != nil {
		domainTask.PurgeAt = mo.Some(*task.PurgeAt)
	}

	if task.NotifyAt != nil {
		domainTask.NotifyAt = task.NotifyAt
	}

	return domainTask, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domaintask.Task, error) {
	task, err := r.queries.GetTaskByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domaintask.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	domainTask, err := toDomainTask(&task)
	if err != nil {
		return nil, fmt.Errorf("convert task: %w", err)
	}

	return domainTask, nil
}

func (r *TaskRepository) Create(ctx context.Context, task *domaintask.Task) error {
	var groupID *uuid.UUID
	if value, ok := task.GroupID.Get(); ok {
		groupID = &value
	}

	var location *geom.Point
	if value, ok := task.Location.Get(); ok {
		location = newPoint(value.Latitude, value.Longitude)
	}

	args := sqlcgen.CreateTaskParams{
		ID:              task.ID,
		OwnerID:         task.OwnerID,
		GroupID:         groupID,
		Title:           task.Title,
		Description:     task.Description,
		Column6:         location,
		IsFavorite:      task.IsFavorite,
		Priority:        task.Priority.String(),
		Icon:            task.Icon.String(),
		Status:          task.Status.String(),
		Deadline:        fromDomainTime(task.Deadline.OrElse(time.Time{})),
		NotifyAt:        task.NotifyAt,
		TitleLang:       r.detector.Detect(task.Title),
		DescriptionLang: r.detector.Detect(task.Description),
	}

	err := r.queries.CreateTask(ctx, args)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *TaskRepository) Save(ctx context.Context, task *domaintask.Task) error {
	var geoPoint *geom.Point
	if value, ok := task.Location.Get(); ok {
		geoPoint = newPoint(value.Latitude, value.Longitude)
	}

	args := sqlcgen.SaveTaskParams{
		ID:              task.ID,
		OwnerID:         task.OwnerID,
		Title:           task.Title,
		Description:     task.Description,
		Location:        geoPoint,
		IsFavorite:      task.IsFavorite,
		Priority:        task.Priority.String(),
		Icon:            task.Icon.String(),
		Status:          task.Status.String(),
		Deadline:        fromDomainTime(task.CreatedAt),
		NotifyAt:        task.NotifyAt,
		CompletedAt:     fromDomainTime(task.CompletedAt.OrElse(time.Time{})),
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
		PurgeAt:         fromDomainTime(task.PurgeAt.OrElse(time.Time{})),
		TitleLang:       r.detector.Detect(task.Title),
		DescriptionLang: r.detector.Detect(task.Description),
	}

	err := r.queries.SaveTask(ctx, args)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}

	return nil
}

func (r *TaskRepository) SoftDelete(ctx context.Context, task *domaintask.Task) error {
	args := sqlcgen.SoftDeleteTaskByIDParams{
		ID:        task.ID,
		PurgeAt:   fromDomainTime(task.PurgeAt.MustGet()),
		UpdatedAt: task.UpdatedAt,
	}

	err := r.queries.SoftDeleteTaskByID(ctx, args)
	if err != nil {
		return fmt.Errorf("soft delete: %w", err)
	}

	return nil
}

func (r *TaskRepository) HardDeleteByUserID(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteTasksByOwnerID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete tasks: %w", err)
	}

	return nil
}

func fromDomainFindDeletedFilter(filter domaintask.FindDeletedFilter) (string, []any, error) {
	query := sqld.Select("id, purge_at").
		From("pgtasks.task").
		Where("purge_at < ?", filter.Cursor.Key).
		PlaceholderFormat(sqld.Dollar)

	if ownerID, ok := filter.OwnerID.Get(); ok {
		query = query.Where("owner_id = ?", ownerID)
	}

	if groupID, ok := filter.GroupID.Get(); ok {
		query = query.Where("group_id = ?", groupID)
	}

	query = query.OrderBy("purge_at DESC").
		Limit(filter.Cursor.Limit)

	return query.ToSql()
}

func (r *TaskRepository) FindDeleted(ctx context.Context, filter domaintask.FindDeletedFilter) ([]uuid.UUID, time.Time, error) {
	query, args, err := fromDomainFindDeletedFilter(filter)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("from domain filter: %w", err)
	}

	rows, err := r.client.Query(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, time.Time{}, domaintask.ErrNotFound
	}

	if err != nil {
		return nil, time.Time{}, fmt.Errorf("query rows: %w", err)
	}

	var ids []uuid.UUID
	var purgeAt time.Time
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id, &purgeAt); err != nil {
			return nil, time.Time{}, fmt.Errorf("scan row: %w", err)
		}

		ids = append(ids, id)
	}

	if rows.Err() != nil {
		return nil, time.Time{}, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(ids) == 0 {
		return nil, time.Time{}, domaintask.ErrNotFound
	}

	return ids, purgeAt, nil
}

func NewTaskRepository(pool *pgxpool.Pool, detector LanguageDetector) (*TaskRepository, error) {
	if pool == nil {
		return nil, errors.New("pool is missing")
	}

	return &TaskRepository{
		client:   pool,
		queries:  sqlcgen.New(pool),
		detector: detector,
	}, nil
}

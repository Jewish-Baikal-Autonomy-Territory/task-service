package mongo

import (
	"context"
	"errors"
	"fmt"
	"task-service/internal/domain/task"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoTaskRepository struct {
	client         TaskClient
	taskCollection string
}

func fromDomainUUID(id uuid.UUID) bson.Binary {
	return bson.Binary{Subtype: 0x04, Data: id[:]}
}

func fromDomainFilter(filter *task.Filter) (bson.M, error) {
	if filter == nil {
		return nil, task.ErrEmptyFilter
	}
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validate filter: %w", err)
	}

	mongoFilter := bson.M{}
	if value, ok := filter.OwnerID.Get(); ok {
		mongoFilter["owner_id"] = bson.Binary{Subtype: 0x04, Data: value[:]}
	}
	if value, ok := filter.GroupID.Get(); ok {
		mongoFilter["group_id"] = bson.Binary{Subtype: 0x04, Data: value[:]}
	}
	if value, ok := filter.IsFavorite.Get(); ok {
		mongoFilter["is_favorite"] = value
	}
	if value, ok := filter.Status.Get(); ok {
		mongoFilter["status"] = value.String()
	}
	if value, ok := filter.Priority.Get(); ok {
		mongoFilter["priority"] = value.String()
	}
	if value, ok := filter.Area.Get(); ok {
		mongoFilter["location"] = bson.D{{
			Key: "$near", Value: bson.D{
				{Key: "$geometry", Value: fromDomainGeoPoint(value.Point)},
				{Key: "$maxDistance", Value: value.Radius},
			},
		},
		}
	}

	return mongoFilter, nil
}

func (mr *mongoTaskRepository) Create(ctx context.Context, task *task.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	doc := fromDomainTask(task)
	_, err := mr.client.
		Database().
		Collection(mr.taskCollection).
		InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (mr *mongoTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	doc := taskDocument{}
	err := mr.client.
		Database().
		Collection(mr.taskCollection).
		FindOne(ctx, bson.M{"_id": bson.Binary{Subtype: 0x04, Data: id[:]}}).
		Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find task by id: %w", err)
	}
	return doc.toDomain(), nil
}

func (mr *mongoTaskRepository) FindAll(ctx context.Context, filter *task.Filter) ([]*task.Task, error) {
	mongoFilter, err := fromDomainFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("build filter: %w", err)
	}

	cursor, err := mr.client.
		Database().
		Collection(mr.taskCollection).
		Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("find tasks: %w", err)
	}
	defer cursor.Close(ctx) //nolint:errcheck

	docs := new([]*taskDocument)
	if err := cursor.All(ctx, docs); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}

	if len(*docs) == 0 {
		return nil, errors.New("tasks not found")
	}

	tasks := make([]*task.Task, len(*docs))
	for i, doc := range *docs {
		tasks[i] = doc.toDomain()
	}
	return tasks, nil
}

func (mr *mongoTaskRepository) FindDeleted(ctx context.Context, filter *task.Filter) ([]*task.Task, error) {
	mongoFilter, err := fromDomainFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("build filter: %w", err)
	}
	mongoFilter["purge_at"] = bson.M{"$ne": nil}

	cursor, err := mr.client.
		Database().
		Collection(mr.taskCollection).
		Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("find deleted tasks: %w", err)
	}
	defer cursor.Close(ctx) //nolint:errcheck

	docs := new([]*taskDocument)
	if err := cursor.All(ctx, docs); err != nil {
		return nil, fmt.Errorf("decode deleted tasks: %w", err)
	}

	if len(*docs) == 0 {
		return nil, errors.New("deleted tasks not found")
	}

	tasks := make([]*task.Task, len(*docs))
	for i, doc := range *docs {
		tasks[i] = doc.toDomain()
	}
	return tasks, nil
}

func (mr *mongoTaskRepository) Update(ctx context.Context, task *task.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}

	doc := fromDomainTask(task)
	result, err := mr.client.
		Database().
		Collection(mr.taskCollection).
		UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (mr *mongoTaskRepository) Close(ctx context.Context) error {
	if err := mr.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect client: %w", err)
	}
	return nil
}

type TaskRepositoryOpts struct {
	Collection string
}

func NewTaskRepository(client TaskClient, opts TaskRepositoryOpts) task.Repository {
	return &mongoTaskRepository{
		client:         client,
		taskCollection: opts.Collection,
	}
}

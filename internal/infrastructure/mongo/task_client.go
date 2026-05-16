package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type TaskClient interface {
	Database() *mongo.Database
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error
}

type mongoTaskClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func (mt *mongoTaskClient) Database() *mongo.Database {
	return mt.db
}

func (mt *mongoTaskClient) Disconnect(ctx context.Context) error {
	return mt.client.Disconnect(ctx)
}

func (mt *mongoTaskClient) Ping(ctx context.Context) error {
	return mt.client.Ping(ctx, readpref.PrimaryPreferred())
}

type taskClientOpts struct {
	connectionString string
	database         string
}

func newTaskClient(opts taskClientOpts) (TaskClient, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(opts.connectionString))
	if err != nil {
		return nil, err
	}
	return &mongoTaskClient{
		client: client,
		db:     client.Database(opts.database),
	}, nil
}

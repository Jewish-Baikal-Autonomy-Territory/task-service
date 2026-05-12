package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Test_mongoTaskClient_Database(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	type fields struct {
		client *mongo.Client
		db     *mongo.Database
	}
	tests := []struct {
		name   string
		fields fields
		want   *mongo.Database
	}{
		{
			name: "valid database",
			fields: fields{
				client: client,
				db:     client.Database("task"),
			},
			want: client.Database("task"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mongoTaskClient{
				client: tt.fields.client,
				db:     tt.fields.db,
			}
			got := mt.Database()
			if got == nil && tt.want != nil || got != nil && tt.want == nil {
				t.Errorf("mongoTaskClient.Database() = %v, want %v", got, tt.want)
			}
			if got.Name() != tt.want.Name() {
				t.Errorf("mongoTaskClient.Database().Name() = %v, want %v", got.Name(), tt.want.Name())
			}
		})
	}
}

func Test_mongoTaskClient_Ping(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type fields struct {
		client *mongo.Client
		db     *mongo.Database
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid client",
			fields: fields{
				client: client,
				db:     client.Database("task"),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "expired context",
			fields: fields{
				client: client,
				db:     client.Database("task"),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mongoTaskClient{
				client: tt.fields.client,
				db:     tt.fields.db,
			}
			if err := mt.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newTaskClient(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	type args struct {
		opts taskClientOpts
	}
	tests := []struct {
		name    string
		args    args
		want    taskClient
		wantErr bool
	}{
		{
			name: "valid connection string",
			args: args{
				opts: taskClientOpts{
					connectionString: "mongodb://localhost:27017",
					database:         "test",
				},
			},
			want: &mongoTaskClient{
				client: client,
				db:     client.Database("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTaskClient(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTaskClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil || got != nil && tt.want == nil {
				t.Errorf("newTaskClient() = %v, want %v", got, tt.want)
			}
			if got != nil && got.Database().Name() != tt.want.Database().Name() {
				t.Errorf("newTaskClient() = %v, want %v", got.Database().Name(), tt.want.Database().Name())
			}
		})
	}
}

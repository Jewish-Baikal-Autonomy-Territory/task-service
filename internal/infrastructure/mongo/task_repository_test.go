package mongo

import (
	"context"
	"errors"
	"log"
	"path"
	"reflect"
	domaintask "task-service/internal/domain/task"
	"task-service/internal/infrastructure/mongo/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	testMongoTaskRepositoryConfig = struct {
		ContainerImage string
		Username       string
		Password       string
		Database       string
		Collection     string
		TestDuration   time.Duration
	}{
		ContainerImage: "mongo:8.0.21",
		Username:       "task-test",
		Password:       "task-test",
		Database:       "test",
		Collection:     "test",
		TestDuration:   10 * time.Second,
	}
)

var testInvalidMongoClient = func() *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}()

type testMongoTaskRepositoryTestSuite struct {
	suite.Suite

	container      *mongodb.MongoDBContainer
	taskRepository mongoTaskRepository
	ctx            context.Context
}

func (ts *testMongoTaskRepositoryTestSuite) SetupSuite() {
	ts.ctx = context.Background()

	container, err := mongodb.Run(
		ts.ctx,
		testMongoTaskRepositoryConfig.ContainerImage,
		testcontainers.WithEnv(map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root-" + testMongoTaskRepositoryConfig.Username,
			"MONGO_INITDB_ROOT_PASSWORD": "root-" + testMongoTaskRepositoryConfig.Password,
			"MONGO_INITDB_DATABASE":      "root-" + testMongoTaskRepositoryConfig.Database,
		}),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      path.Join(".", "testdata", "mongo-init.js"),
			ContainerFilePath: "/docker-entrypoint-initdb.d/mongo-init.js",
			FileMode:          0755,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	ts.container = container
	connString, err := ts.container.ConnectionString(ts.ctx)
	if err != nil {
		log.Fatal(err)
	}
	client, err := newTaskClient(taskClientOpts{
		connectionString: connString,
		database:         testMongoTaskRepositoryConfig.Database,
	})
	if err != nil {
		log.Fatal(err)
	}
	taskRepository := NewTaskRepository(client, TaskRepositoryOpts{
		Collection: testMongoTaskRepositoryConfig.Collection,
	})
	ts.taskRepository = *taskRepository.(*mongoTaskRepository)
}

func (ts *testMongoTaskRepositoryTestSuite) TearDownSuite() {
	if err := ts.taskRepository.Close(ts.ctx); err != nil {
		log.Println(err)
	}
	if err := ts.container.Terminate(ts.ctx); err != nil {
		log.Fatal(err)
	}
}

func (ts *testMongoTaskRepositoryTestSuite) TearDownTest() {
	_, err := ts.taskRepository.client.
		Database().
		Collection(ts.taskRepository.taskCollection).
		DeleteMany(ts.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
}

func (ts *testMongoTaskRepositoryTestSuite) TestCreateTask() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.NoError(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestCreateDuplicateTask() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityMedium)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Error(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestNilCreate() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	err := ts.taskRepository.Create(ctx, nil)
	ts.Error(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindTaskByID() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityHigh)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTask, err := ts.taskRepository.FindByID(ctx, task.ID)

	ts.Require().NoError(err)
	ts.Require().NotNil(foundTask)
	ts.Equal(task.ID, foundTask.ID)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindByID() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := ts.taskRepository.FindByID(ctx, uuid.New())
	ts.Error(err)
	ts.Nil(task)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindAllWithOwnerIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{OwnerID: mo.Some(task.OwnerID)})
	ts.Require().NoError(err)
	ts.Require().NotNil(foundTasks)
	ts.Require().Equal(1, len(foundTasks))
	ts.Equal(foundTasks[0].ID, task.ID)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindAllWithIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{GroupID: mo.Some(uuid.New())})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindAllWithGroupIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	task.GroupID = mo.Some[uuid.UUID](uuid.New())
	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{GroupID: task.GroupID})
	ts.Require().NoError(err)
	ts.Require().NotNil(foundTasks)
	ts.Require().Equal(1, len(foundTasks))
	ts.Equal(foundTasks[0].ID, task.ID)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindAllWithIsFavoriteOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{IsFavorite: mo.Some(true)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindAllWithIsFavoriteOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	task.IsFavorite = true

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{IsFavorite: mo.Some(true)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindAllWithStatusOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{Status: mo.Some(domaintask.StatusPending)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindAllWithStatusOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityHigh)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{Status: mo.Some(task.Status)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindAllWithPriorityOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{Priority: mo.Some(domaintask.PriorityLow)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindAllWithPriorityOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindAll(ctx, &domaintask.Filter{Priority: mo.Some(task.Priority)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestNilFilterFindAll() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	tasks, err := ts.taskRepository.FindAll(ctx, nil)
	ts.Error(err)
	ts.Nil(tasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindDeletedWithOwnerIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{OwnerID: mo.Some(uuid.New())})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindDeletedWithOwnerIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = task.SoftDelete()
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{OwnerID: mo.Some(task.OwnerID)})
	ts.Require().NoError(err)
	ts.Require().NotNil(foundTasks)
	ts.Require().Equal(1, len(foundTasks))
	ts.Equal(foundTasks[0].ID, task.ID)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindDeletedWithGroupIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{GroupID: mo.Some(uuid.New())})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindDeletedWithGroupIDOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	task.GroupID = mo.Some[uuid.UUID](uuid.New())

	err = task.SoftDelete()
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{GroupID: task.GroupID})
	ts.Require().NoError(err)
	ts.Require().NotNil(foundTasks)
	ts.Require().Equal(1, len(foundTasks))
	ts.Equal(foundTasks[0].ID, task.ID)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindDeletedWithIsFavoriteOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{IsFavorite: mo.Some(true)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindDeletedWithIsFavoriteOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = task.SoftDelete()
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{IsFavorite: mo.Some(task.IsFavorite)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindDeletedWithStatusOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{Status: mo.Some(domaintask.StatusPending)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindDeletedWithStatusOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityHigh)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = task.SoftDelete()
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{Status: mo.Some(task.Status)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingFindDeletedWithPriorityOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{Priority: mo.Some(domaintask.PriorityLow)})
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestFindDeletedWithPriorityOnly() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = task.SoftDelete()
	ts.Require().NoError(err)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, &domaintask.Filter{Priority: mo.Some(task.Priority)})
	ts.Require().Error(err)
	ts.Require().Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestNilFilterFindDeleted() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	foundTasks, err := ts.taskRepository.FindDeleted(ctx, nil)
	ts.Error(err)
	ts.Nil(foundTasks)
}

func (ts *testMongoTaskRepositoryTestSuite) TestMissingUpdate() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Update(ctx, task)
	ts.Error(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestInvalidUpdate() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	task.ID = uuid.New()
	err = ts.taskRepository.Update(ctx, task)
	ts.Error(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestUpdate() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	task, err := domaintask.NewTask(uuid.New(), "title", "description", domaintask.PriorityLow)
	ts.Require().NoError(err)
	ts.Require().NotNil(task)

	err = ts.taskRepository.Create(ctx, task)
	ts.Require().NoError(err)

	task.Title = "New Title"
	err = ts.taskRepository.Update(ctx, task)
	ts.NoError(err)
}

func (ts *testMongoTaskRepositoryTestSuite) TestNilUpdate() {
	ctx, cancel := context.WithTimeout(ts.ctx, testMongoTaskRepositoryConfig.TestDuration)
	defer cancel()

	err := ts.taskRepository.Update(ctx, nil)
	ts.Error(err)
}

func TestMongoTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(testMongoTaskRepositoryTestSuite))
}

func Test_fromDomainFilter(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
		testGeoFilter = domaintask.GeoFilter{
			Point: domaintask.GeoPoint{
				Latitude:  67.0,
				Longitude: 52.0,
			},
			Radius: 1.0,
		}
	)
	type args struct {
		filter *domaintask.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    bson.M
		wantErr bool
	}{
		{
			name: "nil filter",
			args: args{
				filter: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner id only",
			args: args{
				filter: &domaintask.Filter{
					OwnerID: mo.Some(testOwnerUUID),
				},
			},
			want:    bson.M{"owner_id": fromDomainUUID(testOwnerUUID)},
			wantErr: false,
		},
		{
			name: "group id only",
			args: args{
				filter: &domaintask.Filter{
					GroupID: mo.Some(testGroupUUID),
				},
			},
			want:    bson.M{"group_id": fromDomainUUID(testGroupUUID)},
			wantErr: false,
		},
		{
			name: "owner and group id only",
			args: args{
				filter: &domaintask.Filter{
					OwnerID: mo.Some(testOwnerUUID),
					GroupID: mo.Some(testGroupUUID),
				},
			},
			want: bson.M{
				"owner_id": fromDomainUUID(testOwnerUUID),
				"group_id": fromDomainUUID(testGroupUUID),
			},
			wantErr: false,
		},
		{
			name: "owner and is favorite",
			args: args{
				filter: &domaintask.Filter{
					OwnerID:    mo.Some(testOwnerUUID),
					IsFavorite: mo.Some(true),
				},
			},
			want: bson.M{
				"owner_id":    fromDomainUUID(testOwnerUUID),
				"is_favorite": true,
			},
			wantErr: false,
		},
		{
			name: "owner and status",
			args: args{
				filter: &domaintask.Filter{
					OwnerID: mo.Some(testOwnerUUID),
					Status:  mo.Some(domaintask.StatusPending),
				},
			},
			want: bson.M{
				"owner_id": fromDomainUUID(testOwnerUUID),
				"status":   domaintask.StatusPending.String(),
			},
			wantErr: false,
		},
		{
			name: "owner and priority",
			args: args{
				filter: &domaintask.Filter{
					OwnerID:  mo.Some(testOwnerUUID),
					Priority: mo.Some(domaintask.PriorityLow),
				},
			},
			want: bson.M{
				"owner_id": fromDomainUUID(testOwnerUUID),
				"priority": domaintask.PriorityLow.String(),
			},
			wantErr: false,
		},
		{
			name: "owner and area",
			args: args{
				filter: &domaintask.Filter{
					OwnerID: mo.Some(testOwnerUUID),
					Area:    mo.Some(testGeoFilter),
				},
			},
			want: bson.M{
				"owner_id": fromDomainUUID(testOwnerUUID),
				"location": bson.D{
					{Key: "$near", Value: bson.D{
						{Key: "$geometry", Value: fromDomainGeoPoint(testGeoFilter.Point)},
						{Key: "$maxDistance", Value: testGeoFilter.Radius},
					},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "group and is favorite",
			args: args{
				filter: &domaintask.Filter{
					GroupID:    mo.Some(testGroupUUID),
					IsFavorite: mo.Some(true),
				},
			},
			want: bson.M{
				"group_id":    fromDomainUUID(testGroupUUID),
				"is_favorite": true,
			},
			wantErr: false,
		},
		{
			name: "group and status",
			args: args{
				filter: &domaintask.Filter{
					GroupID: mo.Some(testGroupUUID),
					Status:  mo.Some(domaintask.StatusPending),
				},
			},
			want: bson.M{
				"group_id": fromDomainUUID(testGroupUUID),
				"status":   domaintask.StatusPending.String(),
			},
			wantErr: false,
		},
		{
			name: "group and priority",
			args: args{
				filter: &domaintask.Filter{
					GroupID:  mo.Some(testGroupUUID),
					Priority: mo.Some(domaintask.PriorityLow),
				},
			},
			want: bson.M{
				"group_id": fromDomainUUID(testGroupUUID),
				"priority": domaintask.PriorityLow.String(),
			},
			wantErr: false,
		},
		{
			name: "group and area",
			args: args{
				filter: &domaintask.Filter{
					GroupID: mo.Some(testGroupUUID),
					Area:    mo.Some(testGeoFilter),
				},
			},
			want: bson.M{
				"group_id": fromDomainUUID(testGroupUUID),
				"location": bson.D{
					{Key: "$near", Value: bson.D{
						{Key: "$geometry", Value: fromDomainGeoPoint(testGeoFilter.Point)},
						{Key: "$maxDistance", Value: testGeoFilter.Radius},
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromDomainFilter(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromDomainFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoTaskRepository_Close(t *testing.T) {
	type fields struct {
		client         TaskClient
		taskCollection string
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
				client: func() TaskClient {
					m := mocks.NewMockTaskClient(t)
					m.EXPECT().Disconnect(mock.Anything).Return(errors.New("error"))
					return m
				}(),
				taskCollection: "task",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mongoTaskRepository{
				client:         tt.fields.client,
				taskCollection: tt.fields.taskCollection,
			}
			if err := mr.Close(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mongoTaskRepository_FindAll(t *testing.T) {
	type fields struct {
		client         TaskClient
		taskCollection string
	}
	type args struct {
		ctx    context.Context
		filter *domaintask.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domaintask.Task
		wantErr bool
	}{
		{
			name: "invalid client",
			fields: fields{
				client: &mongoTaskClient{
					client: testInvalidMongoClient,
					db:     testInvalidMongoClient.Database(testMongoTaskRepositoryConfig.Database),
				},
				taskCollection: testMongoTaskRepositoryConfig.Collection,
			},
			args: args{
				ctx: context.Background(),
				filter: &domaintask.Filter{
					OwnerID: mo.Some(uuid.New()),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mongoTaskRepository{
				client:         tt.fields.client,
				taskCollection: tt.fields.taskCollection,
			}
			got, err := mr.FindAll(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoTaskRepository_FindByID(t *testing.T) {
	type fields struct {
		client         TaskClient
		taskCollection string
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domaintask.Task
		wantErr bool
	}{
		{
			name: "invalid client",
			fields: fields{
				client: &mongoTaskClient{
					client: testInvalidMongoClient,
					db:     testInvalidMongoClient.Database(testMongoTaskRepositoryConfig.Database),
				},
				taskCollection: testMongoTaskRepositoryConfig.Collection,
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mongoTaskRepository{
				client:         tt.fields.client,
				taskCollection: tt.fields.taskCollection,
			}
			got, err := mr.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoTaskRepository_FindDeleted(t *testing.T) {
	type fields struct {
		client         TaskClient
		taskCollection string
	}
	type args struct {
		ctx    context.Context
		filter *domaintask.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domaintask.Task
		wantErr bool
	}{
		{
			name: "invalid client",
			fields: fields{
				client: &mongoTaskClient{
					client: testInvalidMongoClient,
					db:     testInvalidMongoClient.Database(testMongoTaskRepositoryConfig.Database),
				},
				taskCollection: testMongoTaskRepositoryConfig.Collection,
			},
			args: args{
				ctx: context.Background(),
				filter: &domaintask.Filter{
					OwnerID: mo.Some(uuid.New()),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mongoTaskRepository{
				client:         tt.fields.client,
				taskCollection: tt.fields.taskCollection,
			}
			got, err := mr.FindDeleted(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindDeleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindDeleted() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoTaskRepository_Update(t *testing.T) {
	type fields struct {
		client         TaskClient
		taskCollection string
	}
	type args struct {
		ctx  context.Context
		task *domaintask.Task
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
				client: &mongoTaskClient{
					client: testInvalidMongoClient,
					db:     testInvalidMongoClient.Database(testMongoTaskRepositoryConfig.Database),
				},
				taskCollection: testMongoTaskRepositoryConfig.Collection,
			},
			args: args{
				ctx:  context.Background(),
				task: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mongoTaskRepository{
				client:         tt.fields.client,
				taskCollection: tt.fields.taskCollection,
			}
			if err := mr.Update(tt.args.ctx, tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fromDomainUUID(t *testing.T) {
	testUUID := uuid.New()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want bson.Binary
	}{
		{
			name: "valid uuid",
			args: args{
				id: testUUID,
			},
			want: bson.Binary{Subtype: 0x04, Data: testUUID[:]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDomainUUID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

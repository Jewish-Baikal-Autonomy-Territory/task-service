package task

import (
	"context"
	"errors"
	"task-service/internal/domain/task"
	"task-service/internal/domain/task/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	m := mocks.NewMockRepository(t)
	type args struct {
		repository task.Repository
	}
	tests := []struct {
		name    string
		args    args
		want    Service
		wantErr bool
	}{
		{
			name: "invalid repository",
			args: args{
				repository: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid repository",
			args: args{
				repository: m,
			},
			want: &service{
				repository: m,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("NewService() got = %v, want %v", got, tt.want)
			}
			if got != nil && tt.want == nil {
				t.Errorf("NewService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CompleteTask(t *testing.T) {
	var (
		testNotFoundByID           = uuid.New()
		testOwnerID                = uuid.New()
		testMismatchedOwnerID, _   = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testDeletedTask, _         = task.NewTask(testOwnerID, "title", "description", task.PriorityLow)
		testUpdateTask, _          = task.NewTask(testOwnerID, "title", "description", task.PriorityLow)
		testValidTaskCompletion, _ = task.NewTask(testOwnerID, "title", "description", task.PriorityLow)
		expiredCtx, cancel         = context.WithTimeout(context.Background(), 1*time.Nanosecond)
		m                          = mocks.NewMockRepository(t)
	)
	defer cancel()
	time.Sleep(10 * time.Nanosecond)

	require.NotNil(t, testMismatchedOwnerID)
	require.NotNil(t, testDeletedTask)
	require.NotNil(t, testUpdateTask)
	require.NotNil(t, testValidTaskCompletion)
	require.Nil(t, testDeletedTask.SoftDelete())

	m.EXPECT().
		FindByID(mock.Anything, testNotFoundByID).
		Return(nil, errors.New("not found")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testMismatchedOwnerID.ID).
		Return(testMismatchedOwnerID, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testDeletedTask.ID).
		Return(testDeletedTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testUpdateTask.ID).
		Return(testUpdateTask, nil).
		Once()
	m.EXPECT().
		Update(mock.Anything, testUpdateTask).
		Return(errors.New("update error")).
		Once()

	m.EXPECT().
		FindByID(expiredCtx, testValidTaskCompletion.ID).
		Return(nil, errors.New("expired context")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidTaskCompletion.ID).
		Return(testValidTaskCompletion, nil).
		Times(2)
	m.EXPECT().
		Update(expiredCtx, testValidTaskCompletion).
		Return(errors.New("expired context")).
		Once()
	m.EXPECT().
		Update(mock.Anything, testValidTaskCompletion).
		Return(nil).
		Once()

	type fields struct {
		repository task.Repository
	}
	type args struct {
		ctx     context.Context
		command CompleteTaskCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "task not found",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: CompleteTaskCommand{
					ID:      testNotFoundByID,
					OwnerID: testOwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched owner",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: CompleteTaskCommand{
					ID:      testMismatchedOwnerID.ID,
					OwnerID: testOwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "already completed",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: CompleteTaskCommand{
					ID:      testDeletedTask.ID,
					OwnerID: testDeletedTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "update error",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: CompleteTaskCommand{
					ID:      testUpdateTask.ID,
					OwnerID: testUpdateTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[FindByID]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: CompleteTaskCommand{
					ID:      testValidTaskCompletion.ID,
					OwnerID: testValidTaskCompletion.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[Update]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: CompleteTaskCommand{
					ID:      testValidTaskCompletion.ID,
					OwnerID: testValidTaskCompletion.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "valid task completion",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: CompleteTaskCommand{
					ID:      testValidTaskCompletion.ID,
					OwnerID: testValidTaskCompletion.OwnerID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &service{
			repository: tt.fields.repository,
		}
		if err := s.CompleteTask(tt.args.ctx, tt.args.command); (err != nil) != tt.wantErr {
			t.Errorf("%s: CompleteTask() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_service_CreateTask(t *testing.T) {
	var (
		testOwnerID              = uuid.New()
		testValidTaskCreation, _ = task.NewTask(testOwnerID, "title", "description", task.PriorityLow)
		expiredCtx, cancel       = context.WithTimeout(context.Background(), 1*time.Nanosecond)
		m                        = mocks.NewMockRepository(t)
	)
	defer cancel()
	time.Sleep(5 * time.Nanosecond)

	require.NotNil(t, testValidTaskCreation)

	m.EXPECT().
		Create(expiredCtx, mock.Anything).
		Return(errors.New("expired context")).
		Once()

	m.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	type fields struct {
		repository task.Repository
	}
	type args struct {
		ctx     context.Context
		command *CreateTaskCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *task.Task
		wantErr bool
	}{
		{
			name: "missing command",
			fields: fields{
				repository: nil,
			},
			args: args{
				ctx:     context.Background(),
				command: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid task creation",
			fields: fields{
				repository: nil,
			},
			args: args{
				ctx: context.Background(),
				command: &CreateTaskCommand{
					OwnerID:     uuid.Nil,
					Title:       "",
					Description: "",
					IsFavorite:  false,
					Priority:    task.PriorityUnknown,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "expired context",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: &CreateTaskCommand{
					OwnerID:     uuid.New(),
					Title:       "title",
					Description: "description",
					IsFavorite:  true,
					Priority:    task.PriorityHigh,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid task creation",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: &CreateTaskCommand{
					OwnerID:     testValidTaskCreation.OwnerID,
					Title:       testValidTaskCreation.Title,
					Description: testValidTaskCreation.Description,
					IsFavorite:  testValidTaskCreation.IsFavorite,
					Priority:    testValidTaskCreation.Priority,
				},
			},
			want:    testValidTaskCreation,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &service{
			repository: tt.fields.repository,
		}
		got, err := s.CreateTask(tt.args.ctx, tt.args.command)

		require.Conditionf(t,
			func() bool {
				return (err != nil) == tt.wantErr
			},
			"%s: CreateTask() error = %v, wantErr %v", tt.name, err, tt.wantErr,
		)

		require.Conditionf(t,
			func() bool {
				return got == nil && tt.want == nil || got != nil && tt.want != nil
			},
			"%s: CreateTask() got = %v, want = %v", tt.name, got, tt.want,
		)

		if got == nil {
			continue
		}

		require.Equal(t, tt.want.OwnerID, got.OwnerID)
		require.Equal(t, tt.want.GroupID, got.GroupID)
		require.Equal(t, tt.want.Title, got.Title)
		require.Equal(t, tt.want.Description, got.Description)
		require.Equal(t, tt.want.IsFavorite, got.IsFavorite)
		require.Equal(t, tt.want.Priority, got.Priority)
		require.True(t, tt.want.Deadline.IsNone() && got.Deadline.IsNone())
		require.True(t, tt.want.PurgeAt.IsNone() && got.PurgeAt.IsNone())
	}
}

func Test_service_DeleteTask(t *testing.T) {
	var (
		testMissingTaskID            = uuid.New()
		testMismatchedOwnerIDTask, _ = task.NewTask(uuid.New(), "title", "description", task.PriorityMedium)
		testAlreadyDeletedTask, _    = task.NewTask(uuid.New(), "title", "description", task.PriorityHigh)
		testValidExpiredTask, _      = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testValidTask, _             = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		expiredCtx, cancel           = context.WithTimeout(context.Background(), 1*time.Nanosecond)
		m                            = mocks.NewMockRepository(t)
	)
	defer cancel()
	time.Sleep(5 * time.Nanosecond)

	require.NotNil(t, testMismatchedOwnerIDTask)
	require.NotNil(t, testAlreadyDeletedTask)
	require.NotNil(t, testValidExpiredTask)
	require.NotNil(t, testValidTask)

	_ = testAlreadyDeletedTask.SoftDelete()

	m.EXPECT().
		FindByID(mock.Anything, testMissingTaskID).
		Return(nil, errors.New("not found")).
		Once()

	m.EXPECT().
		FindByID(expiredCtx, mock.Anything).
		Return(nil, errors.New("expired context")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testMismatchedOwnerIDTask.ID).
		Return(testMismatchedOwnerIDTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testAlreadyDeletedTask.ID).
		Return(testAlreadyDeletedTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidExpiredTask.ID).
		Return(testValidExpiredTask, nil).
		Once()
	m.EXPECT().
		Update(expiredCtx, mock.Anything).
		Return(errors.New("expired context")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidTask.ID).
		Return(testValidTask, nil).
		Once()
	m.EXPECT().
		Update(mock.Anything, testValidTask).
		Return(nil).
		Once()

	type fields struct {
		repository task.Repository
	}
	type args struct {
		ctx     context.Context
		command DeleteTaskCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "missing task",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: DeleteTaskCommand{
					ID:      testMissingTaskID,
					OwnerID: uuid.Nil,
				},
			},
			wantErr: true,
		},
		{
			name: "context expired[FindByID]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: DeleteTaskCommand{
					ID:      uuid.New(),
					OwnerID: uuid.Nil,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched owner",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: DeleteTaskCommand{
					ID:      testMismatchedOwnerIDTask.ID,
					OwnerID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "already deleted task",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: DeleteTaskCommand{
					ID:      testAlreadyDeletedTask.ID,
					OwnerID: testAlreadyDeletedTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[Update]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: DeleteTaskCommand{
					ID:      testValidExpiredTask.ID,
					OwnerID: testValidExpiredTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "valid task deletion",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: DeleteTaskCommand{
					ID:      testValidTask.ID,
					OwnerID: testValidTask.OwnerID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repository: tt.fields.repository,
			}
			if err := s.DeleteTask(tt.args.ctx, tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_RestoreTask(t *testing.T) {
	var (
		testNotFoundByID             = uuid.New()
		testMismatchedOwnerIDTask, _ = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testNotDeletedTask, _        = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testValidExpiredTask, _      = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testValidTask, _             = task.NewTask(uuid.New(), "title", "description", task.PriorityHigh)
		expiredCtx, cancel           = context.WithTimeout(context.Background(), 1*time.Nanosecond)
		m                            = mocks.NewMockRepository(t)
	)

	defer cancel()
	time.Sleep(5 * time.Nanosecond)

	require.NotNil(t, testMismatchedOwnerIDTask)
	require.NotNil(t, testNotDeletedTask)
	require.NotNil(t, testValidExpiredTask)
	require.NotNil(t, testValidTask)

	_ = testValidExpiredTask.SoftDelete()
	_ = testValidTask.SoftDelete()

	m.EXPECT().
		FindByID(mock.Anything, testNotFoundByID).
		Return(nil, errors.New("not found")).
		Once()

	m.EXPECT().
		FindByID(expiredCtx, mock.Anything).
		Return(nil, errors.New("context expired")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testMismatchedOwnerIDTask.ID).
		Return(testMismatchedOwnerIDTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testNotDeletedTask.ID).
		Return(testNotDeletedTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidExpiredTask.ID).
		Return(testValidExpiredTask, nil).
		Once()
	m.EXPECT().
		Update(expiredCtx, testValidExpiredTask).
		Return(errors.New("context expired")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidTask.ID).
		Return(testValidTask, nil).
		Once()
	m.EXPECT().
		Update(mock.Anything, testValidTask).
		Return(nil).
		Once()

	type fields struct {
		repository task.Repository
	}
	type args struct {
		ctx     context.Context
		command RestoreTaskCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "missing task",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: RestoreTaskCommand{
					ID:      testNotFoundByID,
					OwnerID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "context expired[FindByID]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: RestoreTaskCommand{
					ID:      uuid.New(),
					OwnerID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched owner",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: RestoreTaskCommand{
					ID:      testMismatchedOwnerIDTask.ID,
					OwnerID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "not deleted task",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: RestoreTaskCommand{
					ID:      testNotDeletedTask.ID,
					OwnerID: testNotDeletedTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[Update]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: RestoreTaskCommand{
					ID:      testValidExpiredTask.ID,
					OwnerID: testValidExpiredTask.OwnerID,
				},
			},
			wantErr: true,
		},
		{
			name: "valid task restore",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: RestoreTaskCommand{
					ID:      testValidTask.ID,
					OwnerID: testValidTask.OwnerID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &service{
			repository: tt.fields.repository,
		}
		err := s.RestoreTask(tt.args.ctx, tt.args.command)

		require.Conditionf(t,
			func() bool {
				return (err != nil) == tt.wantErr
			},
			"%s: RestoreTask() error = %v, wantErr = %v", tt.name, err, tt.wantErr,
		)
	}
}

func Test_service_UpdateTask(t *testing.T) {
	var (
		testMissingID              = uuid.New()
		testMismatchedOwnerTask, _ = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testValidExpiredTask, _    = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		testValidTaskUpdate, _     = task.NewTask(uuid.New(), "title", "description", task.PriorityLow)
		expiredCtx, cancel         = context.WithTimeout(context.Background(), 1*time.Nanosecond)
		m                          = mocks.NewMockRepository(t)
	)
	defer cancel()
	time.Sleep(5 * time.Nanosecond)

	require.NotNil(t, testMismatchedOwnerTask)
	require.NotNil(t, testValidExpiredTask)
	require.NotNil(t, testValidTaskUpdate)

	m.EXPECT().
		FindByID(mock.Anything, testMissingID).
		Return(nil, errors.New("not found")).
		Once()

	m.EXPECT().
		FindByID(expiredCtx, mock.Anything).
		Return(nil, errors.New("context expired")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testMismatchedOwnerTask.ID).
		Return(testMismatchedOwnerTask, nil).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidExpiredTask.ID).
		Return(testValidExpiredTask, nil).
		Once()
	m.EXPECT().
		Update(expiredCtx, testValidExpiredTask).
		Return(errors.New("context expired")).
		Once()

	m.EXPECT().
		FindByID(mock.Anything, testValidTaskUpdate.ID).
		Return(testValidTaskUpdate, nil).
		Once()
	m.EXPECT().
		Update(mock.Anything, testValidTaskUpdate).
		Return(nil).
		Once()

	type fields struct {
		repository task.Repository
	}
	type args struct {
		ctx     context.Context
		command *UpdateTaskCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "missing command",
			fields: fields{
				repository: nil,
			},
			args: args{
				ctx:     context.Background(),
				command: nil,
			},
			wantErr: true,
		},
		{
			name: "missing task",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: &UpdateTaskCommand{
					ID: testMissingID,
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[FindByID]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: &UpdateTaskCommand{
					ID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched owner",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: &UpdateTaskCommand{
					ID:      testMismatchedOwnerTask.ID,
					OwnerID: uuid.New(),
				},
			},
			wantErr: true,
		},
		{
			name: "expired context[Update]",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: expiredCtx,
				command: &UpdateTaskCommand{
					ID:          testValidExpiredTask.ID,
					OwnerID:     testValidExpiredTask.OwnerID,
					GroupID:     mo.Some(uuid.New()),
					Title:       mo.Some("another-title"),
					Description: mo.Some("another-description"),
					Priority:    mo.Some(task.PriorityMedium),
					IsFavorite:  mo.Some(true),
				},
			},
			wantErr: true,
		},
		{
			name: "valid task update",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx: context.Background(),
				command: &UpdateTaskCommand{
					ID:          testValidTaskUpdate.ID,
					OwnerID:     testValidTaskUpdate.OwnerID,
					GroupID:     mo.Some(uuid.New()),
					Title:       mo.Some("another-title"),
					Description: mo.Some("another-description"),
					Priority:    mo.Some(task.PriorityHigh),
					IsFavorite:  mo.Some(false),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &service{
			repository: tt.fields.repository,
		}
		err := s.UpdateTask(tt.args.ctx, tt.args.command)
		require.Truef(t,
			(err != nil) == tt.wantErr,
			"%s: UpdateTask() error = %v, wantErr = %v", tt.name, err, tt.wantErr,
		)
	}
}

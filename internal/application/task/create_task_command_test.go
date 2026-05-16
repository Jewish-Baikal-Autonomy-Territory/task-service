package task

import (
	"reflect"
	"task-service/internal/domain/task"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

func TestCreateTaskCommandBuilder_Build(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	tests := []struct {
		name    string
		fields  fields
		want    *CreateTaskCommand
		wantErr bool
	}{
		{
			name:    "empty initialization",
			fields:  fields{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "title only",
			fields: fields{
				title: "title",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityLow,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner and title",
			fields: fields{
				ownerID: testOwnerUUID,
				title:   "title",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner and priority",
			fields: fields{
				ownerID:  testOwnerUUID,
				priority: task.PriorityLow,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "title and priority",
			fields: fields{
				title:    "title",
				priority: task.PriorityLow,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner, title and priority",
			fields: fields{
				ownerID:  testOwnerUUID,
				title:    "title",
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommand{
				OwnerID:  testOwnerUUID,
				Title:    "title",
				Priority: task.PriorityLow,
			},
			wantErr: false,
		},
		{
			name: "owner, title, priority and group",
			fields: fields{
				ownerID:  testOwnerUUID,
				groupID:  mo.Some(testGroupUUID),
				title:    "title",
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommand{
				OwnerID:  testOwnerUUID,
				GroupID:  mo.Some(testGroupUUID),
				Title:    "title",
				Priority: task.PriorityLow,
			},
			wantErr: false,
		},
		{
			name: "owner, title, priority and description",
			fields: fields{
				ownerID:     testOwnerUUID,
				title:       "title",
				description: "description",
				priority:    task.PriorityLow,
			},
			want: &CreateTaskCommand{
				OwnerID:     testOwnerUUID,
				Title:       "title",
				Description: "description",
				Priority:    task.PriorityLow,
			},
			wantErr: false,
		},
		{
			name: "owner, title, priority and is favorite",
			fields: fields{
				ownerID:    testOwnerUUID,
				title:      "title",
				isFavorite: true,
				priority:   task.PriorityLow,
			},
			want: &CreateTaskCommand{
				OwnerID:    testOwnerUUID,
				Title:      "title",
				IsFavorite: true,
				Priority:   task.PriorityLow,
			},
			wantErr: false,
		},
		{
			name: "fully initialized",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			want: &CreateTaskCommand{
				OwnerID:     testOwnerUUID,
				GroupID:     mo.Some(testGroupUUID),
				Title:       "title",
				Description: "description",
				IsFavorite:  true,
				Priority:    task.PriorityLow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			got, err := b.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetDescription(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		description string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				description: "description",
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				description: "description",
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				groupID:     mo.Some(testGroupUUID),
				description: "description",
			},
		},
		{
			name: "title only",
			fields: fields{
				title: "title",
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				title:       "title",
				description: "description",
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				isFavorite:  true,
				description: "description",
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityLow,
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				description: "description",
				priority:    task.PriorityLow,
			},
		},
		{
			name: "initialized without description",
			fields: fields{
				ownerID:    testOwnerUUID,
				groupID:    mo.Some(testGroupUUID),
				title:      "title",
				isFavorite: true,
				priority:   task.PriorityLow,
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
		{
			name: "initialized with description",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "another-description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				description: "description",
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetDescription(tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetGroupID(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				groupID: mo.Some(testGroupUUID),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID: testOwnerUUID,
				groupID: mo.Some(testGroupUUID),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: "title",
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				groupID: mo.Some(testGroupUUID),
				title:   "title",
			},
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				groupID:     mo.Some(testGroupUUID),
				description: "description",
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				groupID:    mo.Some(testGroupUUID),
				isFavorite: true,
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityLow,
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				groupID:  mo.Some(testGroupUUID),
				priority: task.PriorityLow,
			},
		},
		{
			name: "initialized without group",
			fields: fields{
				ownerID:     testOwnerUUID,
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
		{
			name: "initialized with group",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(uuid.New()),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				id: testGroupUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetGroupID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetIsFavorite(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		isFavorite bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization (true)",
			fields: fields{},
			args: args{
				isFavorite: true,
			},
			want: &CreateTaskCommandBuilder{
				isFavorite: true,
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			args: args{
				isFavorite: true,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:    testOwnerUUID,
				isFavorite: true,
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID:    mo.Some(testGroupUUID),
				isFavorite: true,
			},
			args: args{
				isFavorite: false,
			},
			want: &CreateTaskCommandBuilder{
				groupID:    mo.Some(testGroupUUID),
				isFavorite: false,
			},
		},
		{
			name: "title only",
			fields: fields{
				title:      "title",
				isFavorite: true,
			},
			args: args{
				isFavorite: false,
			},
			want: &CreateTaskCommandBuilder{
				title:      "title",
				isFavorite: false,
			},
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
				isFavorite:  true,
			},
			args: args{
				isFavorite: false,
			},
			want: &CreateTaskCommandBuilder{
				description: "description",
				isFavorite:  false,
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityLow,
			},
			args: args{
				isFavorite: false,
			},
			want: &CreateTaskCommandBuilder{
				isFavorite: false,
				priority:   task.PriorityLow,
			},
		},
		{
			name: "initialized without is favorite",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				priority:    task.PriorityLow,
			},
			args: args{
				isFavorite: true,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
		{
			name: "initialized with is favorite",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				isFavorite: true,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetIsFavorite(tt.args.isFavorite); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIsFavorite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetOwnerID(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID: testOwnerUUID,
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID: testOwnerUUID,
				groupID: mo.Some(testGroupUUID),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: "title",
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID: testOwnerUUID,
				title:   "title",
			},
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				description: "description",
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:    testOwnerUUID,
				isFavorite: true,
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityLow,
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:  testOwnerUUID,
				priority: task.PriorityLow,
			},
		},
		{
			name: "initialized without owner",
			fields: fields{
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
		{
			name: "initialized with owner",
			fields: fields{
				ownerID:     uuid.New(),
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetOwnerID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetOwnerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetPriority(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		priority task.Priority
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommandBuilder{
				priority: task.PriorityLow,
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:  testOwnerUUID,
				priority: task.PriorityLow,
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommandBuilder{
				groupID:  mo.Some(testGroupUUID),
				priority: task.PriorityLow,
			},
		},
		{
			name: "title only",
			fields: fields{
				title: "title",
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommandBuilder{
				title:    "title",
				priority: task.PriorityLow,
			},
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &CreateTaskCommandBuilder{
				description: "description",
				priority:    task.PriorityLow,
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			args: args{
				priority: task.PriorityMedium,
			},
			want: &CreateTaskCommandBuilder{
				isFavorite: true,
				priority:   task.PriorityMedium,
			},
		},
		{
			name: "initialized without priority",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
			},
			args: args{
				priority: task.PriorityHigh,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityHigh,
			},
		},
		{
			name: "initialized with priority",
			fields: fields{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityLow,
			},
			args: args{
				priority: task.PriorityMedium,
			},
			want: &CreateTaskCommandBuilder{
				ownerID:     testOwnerUUID,
				groupID:     mo.Some(testGroupUUID),
				title:       "title",
				description: "description",
				isFavorite:  true,
				priority:    task.PriorityMedium,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetPriority(tt.args.priority); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTaskCommandBuilder_SetTitle(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
	)
	type fields struct {
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       string
		description string
		isFavorite  bool
		priority    task.Priority
	}
	type args struct {
		title string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				title: "title",
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerUUID,
			},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				ownerID: testOwnerUUID,
				title:   "title",
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				groupID: mo.Some(testGroupUUID),
				title:   "title",
			},
		},
		{
			name: "description only",
			fields: fields{
				description: "description",
			},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				title:       "title",
				description: "description",
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: true,
			},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				title:      "title",
				isFavorite: true,
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: task.PriorityMedium,
			},
			args: args{
				title: "title",
			},
			want: &CreateTaskCommandBuilder{
				title:    "title",
				priority: task.PriorityMedium,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &CreateTaskCommandBuilder{
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				isFavorite:  tt.fields.isFavorite,
				priority:    tt.fields.priority,
			}
			if got := b.SetTitle(tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

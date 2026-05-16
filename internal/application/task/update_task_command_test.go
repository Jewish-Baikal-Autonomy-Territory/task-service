package task

import (
	"reflect"
	"task-service/internal/domain/task"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

func TestUpdateTaskCommandBuilder_Build(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	tests := []struct {
		name    string
		fields  fields
		want    *UpdateTaskCommand
		wantErr bool
	}{
		{
			name:    "empty initialization",
			fields:  fields{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "group only",
			fields: fields{
				ownerID: testOwnerID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(true),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "id and owner only",
			fields: fields{
				id:      testID,
				ownerID: testOwnerID,
			},
			want: &UpdateTaskCommand{
				ID:      testID,
				OwnerID: testOwnerID,
			},
			wantErr: false,
		},
		{
			name: "id, owner and group only",
			fields: fields{
				id:      testID,
				ownerID: testOwnerID,
				groupID: mo.Some(testGroupID),
			},
			want: &UpdateTaskCommand{
				ID:      testID,
				OwnerID: testOwnerID,
				GroupID: mo.Some(testGroupID),
			},
			wantErr: false,
		},
		{
			name: "id, owner and title only",
			fields: fields{
				id:      testID,
				ownerID: testOwnerID,
				title:   mo.Some("title"),
			},
			want: &UpdateTaskCommand{
				ID:      testID,
				OwnerID: testOwnerID,
				Title:   mo.Some("title"),
			},
			wantErr: false,
		},
		{
			name: "id, owner and description only",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				description: mo.Some("description"),
			},
			want: &UpdateTaskCommand{
				ID:          testID,
				OwnerID:     testOwnerID,
				Description: mo.Some("description"),
			},
			wantErr: false,
		},
		{
			name: "id, owner and is favorite only",
			fields: fields{
				id:         testID,
				ownerID:    testOwnerID,
				isFavorite: mo.Some(true),
			},
			want: &UpdateTaskCommand{
				ID:         testID,
				OwnerID:    testOwnerID,
				IsFavorite: mo.Some(true),
			},
			wantErr: false,
		},
		{
			name: "id, owner and priority only",
			fields: fields{
				id:       testID,
				ownerID:  testOwnerID,
				priority: mo.Some(task.PriorityLow),
			},
			want: &UpdateTaskCommand{
				ID:       testID,
				OwnerID:  testOwnerID,
				Priority: mo.Some(task.PriorityLow),
			},
			wantErr: false,
		},
		{
			name: "full initialization",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			want: &UpdateTaskCommand{
				ID:          testID,
				OwnerID:     testOwnerID,
				GroupID:     mo.Some(testGroupID),
				Title:       mo.Some("title"),
				Description: mo.Some("description"),
				IsFavorite:  mo.Some(true),
				Priority:    mo.Some(task.PriorityLow),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
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

func TestUpdateTaskCommandBuilder_SetDescription(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		description string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				description: mo.Some("description"),
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				description: mo.Some("description"),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:     testOwnerID,
				description: mo.Some("description"),
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				groupID:     mo.Some(testGroupID),
				description: mo.Some("description"),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				title:       mo.Some("title"),
				description: mo.Some("description"),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(false),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				description: mo.Some("description"),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without description",
			fields: fields{
				id:         testID,
				ownerID:    testOwnerID,
				groupID:    mo.Some(testGroupID),
				title:      mo.Some("title"),
				isFavorite: mo.Some(true),
				priority:   mo.Some(task.PriorityLow),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized with description",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("another-description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				description: "description",
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetDescription(tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetGroupID(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				groupID: mo.Some(testGroupID),
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				id:      testID,
				groupID: mo.Some(testGroupID),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID: testOwnerID,
				groupID: mo.Some(testGroupID),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				groupID: mo.Some(testGroupID),
				title:   mo.Some("title"),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				groupID:     mo.Some(testGroupID),
				description: mo.Some("description"),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(false),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				groupID:    mo.Some(testGroupID),
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				groupID:  mo.Some(testGroupID),
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without group",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized with group",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(uuid.New()),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetGroupID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetID(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id: testID,
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:      testID,
				ownerID: testOwnerID,
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:      testID,
				groupID: mo.Some(testGroupID),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				id: testGroupID,
			},
			want: &UpdateTaskCommandBuilder{
				id:    testGroupID,
				title: mo.Some("title"),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				description: mo.Some("description"),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(false),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:         testID,
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:       testID,
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without id",
			fields: fields{
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized with group",
			fields: fields{
				id:          uuid.New(),
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
			args: args{
				id: testID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetIsFavorite(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		isFavorite bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				isFavorite: false,
			},
			want: &UpdateTaskCommandBuilder{
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				isFavorite: true,
			},
			want: &UpdateTaskCommandBuilder{
				id:         testID,
				isFavorite: mo.Some(true),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				isFavorite: false,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:    testOwnerID,
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				isFavorite: false,
			},
			want: &UpdateTaskCommandBuilder{
				groupID:    mo.Some(testGroupID),
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				isFavorite: false,
			},
			want: &UpdateTaskCommandBuilder{
				title:      mo.Some("title"),
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				isFavorite: true,
			},
			want: &UpdateTaskCommandBuilder{
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			args: args{
				isFavorite: true,
			},
			want: &UpdateTaskCommandBuilder{
				isFavorite: mo.Some(true),
				priority:   mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without is favorite",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				priority:    mo.Some(task.PriorityMedium),
			},
			args: args{
				isFavorite: false,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityMedium),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetIsFavorite(tt.args.isFavorite); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIsFavorite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetOwnerID(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID: testOwnerID,
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				id:      testID,
				ownerID: testOwnerID,
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID: testOwnerID,
				groupID: mo.Some(testGroupID),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID: testOwnerID,
				title:   mo.Some("title"),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:     testOwnerID,
				description: mo.Some("description"),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(false),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:    testOwnerID,
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityMedium),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:  testOwnerID,
				priority: mo.Some(task.PriorityMedium),
			},
		},
		{
			name: "initialized without owner",
			fields: fields{
				id:          testID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized with owner",
			fields: fields{
				id:          testID,
				ownerID:     uuid.New(),
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
			args: args{
				id: testOwnerID,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetOwnerID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetOwnerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetPriority(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		priority task.Priority
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				priority: task.PriorityLow,
			},
			want: &UpdateTaskCommandBuilder{
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &UpdateTaskCommandBuilder{
				id:       testID,
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &UpdateTaskCommandBuilder{
				ownerID:  testOwnerID,
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &UpdateTaskCommandBuilder{
				groupID:  mo.Some(testGroupID),
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "title only",
			fields: fields{
				title: mo.Some("title"),
			},
			args: args{
				priority: task.PriorityMedium,
			},
			want: &UpdateTaskCommandBuilder{
				title:    mo.Some("title"),
				priority: mo.Some(task.PriorityMedium),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				priority: task.PriorityHigh,
			},
			want: &UpdateTaskCommandBuilder{
				description: mo.Some("description"),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(true),
			},
			args: args{
				priority: task.PriorityLow,
			},
			want: &UpdateTaskCommandBuilder{
				isFavorite: mo.Some(true),
				priority:   mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without priority",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
			},
			args: args{
				priority: task.PriorityHigh,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
		{
			name: "initialized with priority",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				priority: task.PriorityHigh,
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetPriority(tt.args.priority); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskCommandBuilder_SetTitle(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
		testGroupID = uuid.New()
	)
	type fields struct {
		id          uuid.UUID
		ownerID     uuid.UUID
		groupID     mo.Option[uuid.UUID]
		title       mo.Option[string]
		description mo.Option[string]
		priority    mo.Option[task.Priority]
		isFavorite  mo.Option[bool]
	}
	type args struct {
		title string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UpdateTaskCommandBuilder
	}{
		{
			name:   "empty initialization",
			fields: fields{},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				title: mo.Some("title"),
			},
		},
		{
			name: "id only",
			fields: fields{
				id: testID,
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				id:    testID,
				title: mo.Some("title"),
			},
		},
		{
			name: "owner only",
			fields: fields{
				ownerID: testOwnerID,
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				ownerID: testOwnerID,
				title:   mo.Some("title"),
			},
		},
		{
			name: "group only",
			fields: fields{
				groupID: mo.Some(testGroupID),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				groupID: mo.Some(testGroupID),
				title:   mo.Some("title"),
			},
		},
		{
			name: "description only",
			fields: fields{
				description: mo.Some("description"),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				title:       mo.Some("title"),
				description: mo.Some("description"),
			},
		},
		{
			name: "is favorite only",
			fields: fields{
				isFavorite: mo.Some(true),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				title:      mo.Some("title"),
				isFavorite: mo.Some(true),
			},
		},
		{
			name: "priority only",
			fields: fields{
				priority: mo.Some(task.PriorityLow),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				title:    mo.Some("title"),
				priority: mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized without title",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(true),
				priority:    mo.Some(task.PriorityLow),
			},
		},
		{
			name: "initialized with title",
			fields: fields{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("another-title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
			args: args{
				title: "title",
			},
			want: &UpdateTaskCommandBuilder{
				id:          testID,
				ownerID:     testOwnerID,
				groupID:     mo.Some(testGroupID),
				title:       mo.Some("title"),
				description: mo.Some("description"),
				isFavorite:  mo.Some(false),
				priority:    mo.Some(task.PriorityHigh),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &UpdateTaskCommandBuilder{
				id:          tt.fields.id,
				ownerID:     tt.fields.ownerID,
				groupID:     tt.fields.groupID,
				title:       tt.fields.title,
				description: tt.fields.description,
				priority:    tt.fields.priority,
				isFavorite:  tt.fields.isFavorite,
			}
			if got := b.SetTitle(tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

func TestNewTask(t *testing.T) {
	type args struct {
		ownerID     uuid.UUID
		title       string
		description string
		priority    Priority
	}
	tests := []struct {
		name    string
		args    args
		want    *Task
		wantErr bool
	}{
		{
			name: "invalid title",
			args: args{
				ownerID:     uuid.New(),
				title:       "",
				description: "",
				priority:    PriorityUnknown,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid description",
			args: args{
				ownerID:     uuid.New(),
				title:       "title",
				description: "",
				priority:    PriorityUnknown,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid priority",
			args: args{
				ownerID:     uuid.New(),
				title:       "title",
				description: "description",
				priority:    PriorityUnknown,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid user",
			args: args{
				ownerID:     uuid.New(),
				title:       "title",
				description: "description",
				priority:    PriorityLow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTask(tt.args.ownerID, tt.args.title, tt.args.description, tt.args.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (tt.want != nil) &&
				(got.OwnerID != tt.want.OwnerID ||
					got.Title != tt.want.Title ||
					got.Description != tt.want.Description ||
					got.Priority != tt.want.Priority) {
				t.Errorf("NewTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_CanRestore(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "invalid restore state",
			fields: fields{
				ID:          uuid.New(),
				OwnerID:     uuid.New(),
				Title:       "title",
				Description: "description",
				Priority:    PriorityLow,
			},
			want: false,
		},
		{
			name: "expired restore state",
			fields: fields{
				ID:          uuid.New(),
				OwnerID:     uuid.New(),
				Title:       "title",
				Description: "description",
				Priority:    PriorityLow,
				PurgeAt:     mo.Some[time.Time](time.Now()),
			},
		},
		{
			name: "valid restore",
			fields: fields{
				ID:          uuid.New(),
				OwnerID:     uuid.New(),
				Title:       "title",
				Description: "description",
				Priority:    PriorityLow,
				PurgeAt:     mo.Some[time.Time](time.Now().Add(restoreWindow)),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if got := t.CanRestore(); got != tt.want {
				t1.Errorf("CanRestore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Complete(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid complete state",
			fields: fields{
				PurgeAt: mo.Some[time.Time](time.Now()),
			},
			wantErr: true,
		},
		{
			name:    "valid complete state",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if err := t.Complete(); (err != nil) != tt.wantErr {
				t1.Errorf("Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_IsCompleted(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "invalid complete state",
			fields: fields{},
			want:   false,
		},
		{
			name: "valid complete state",
			fields: fields{
				CompletedAt: mo.Some[time.Time](time.Now()),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if got := t.IsCompleted(); got != tt.want {
				t1.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsDeleted(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "invalid deleted state",
			fields: fields{},
			want:   false,
		},
		{
			name: "valid deleted state",
			fields: fields{
				PurgeAt: mo.Some[time.Time](time.Now()),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if got := t.IsDeleted(); got != tt.want {
				t1.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Restore(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "not deleted task state",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "expired restore task state",
			fields: fields{
				PurgeAt: mo.Some[time.Time](time.Time{}),
			},
			wantErr: true,
		},
		{
			name: "valid restore task state",
			fields: fields{
				PurgeAt: mo.Some[time.Time](time.Now().Add(restoreWindow)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if err := t.Restore(); (err != nil) != tt.wantErr {
				t1.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_SoftDelete(t1 *testing.T) {
	type fields struct {
		ID          uuid.UUID
		OwnerID     uuid.UUID
		GroupID     mo.Option[uuid.UUID]
		Title       string
		Description string
		Location    mo.Option[GeoPoint]
		IsFavorite  bool
		Priority    Priority
		Status      Status
		CreatedAt   time.Time
		UpdatedAt   time.Time
		CompletedAt mo.Option[time.Time]
		Deadline    mo.Option[time.Time]
		PurgeAt     mo.Option[time.Time]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "already deleted task state",
			fields: fields{
				PurgeAt: mo.Some[time.Time](time.Now()),
			},
			wantErr: true,
		},
		{
			name:    "valid soft delete state",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				OwnerID:     tt.fields.OwnerID,
				GroupID:     tt.fields.GroupID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Location:    tt.fields.Location,
				IsFavorite:  tt.fields.IsFavorite,
				Priority:    tt.fields.Priority,
				Status:      tt.fields.Status,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				CompletedAt: tt.fields.CompletedAt,
				Deadline:    tt.fields.Deadline,
				PurgeAt:     tt.fields.PurgeAt,
			}
			if err := t.SoftDelete(); (err != nil) != tt.wantErr {
				t1.Errorf("SoftDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

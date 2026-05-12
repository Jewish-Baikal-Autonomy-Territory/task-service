package task

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

func TestFilterBuilder_Build(t *testing.T) {
	var (
		testOwnerUUID = uuid.New()
		testGroupUUID = uuid.New()
		testGeoFilter = GeoFilter{
			Point: GeoPoint{
				Latitude:  67.0,
				Longitude: 52.0,
			},
			Radius: 5.0,
		}
	)
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Filter
		wantErr bool
	}{
		{
			name:    "empty filter",
			fields:  fields{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only owner id",
			fields: fields{
				ownerID: mo.Some(testOwnerUUID),
			},
			want: &Filter{
				OwnerID: mo.Some(testOwnerUUID),
			},
			wantErr: false,
		},
		{
			name: "only group id",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
			},
			want: &Filter{
				GroupID: mo.Some(testGroupUUID),
			},
			wantErr: false,
		},
		{
			name: "only is favorite",
			fields: fields{
				isFavorite: mo.Some(true),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only status",
			fields: fields{
				status: mo.Some(StatusPending),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only priority",
			fields: fields{
				priority: mo.Some(PriorityLow),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only area",
			fields: fields{
				area: mo.Some(testGeoFilter),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only owner and group id",
			fields: fields{
				ownerID: mo.Some(testOwnerUUID),
				groupID: mo.Some(testGroupUUID),
			},
			want: &Filter{
				OwnerID: mo.Some(testOwnerUUID),
				GroupID: mo.Some(testGroupUUID),
			},
			wantErr: false,
		},
		{
			name: "only owner and is favorite",
			fields: fields{
				ownerID:    mo.Some(testOwnerUUID),
				isFavorite: mo.Some(true),
			},
			want: &Filter{
				OwnerID:    mo.Some(testOwnerUUID),
				IsFavorite: mo.Some(true),
			},
			wantErr: false,
		},
		{
			name: "only owner and status",
			fields: fields{
				ownerID: mo.Some(testOwnerUUID),
				status:  mo.Some(StatusPending),
			},
			want: &Filter{
				OwnerID: mo.Some(testOwnerUUID),
				Status:  mo.Some(StatusPending),
			},
			wantErr: false,
		},
		{
			name: "only owner and priority",
			fields: fields{
				ownerID:  mo.Some(testOwnerUUID),
				priority: mo.Some(PriorityLow),
			},
			want: &Filter{
				OwnerID:  mo.Some(testOwnerUUID),
				Priority: mo.Some(PriorityLow),
			},
			wantErr: false,
		},
		{
			name: "only owner and area",
			fields: fields{
				ownerID: mo.Some(testOwnerUUID),
				area:    mo.Some(testGeoFilter),
			},
			want: &Filter{
				OwnerID: mo.Some(testOwnerUUID),
				Area:    mo.Some(testGeoFilter),
			},
			wantErr: false,
		},
		{
			name: "only group and is favorite",
			fields: fields{
				groupID:    mo.Some(testGroupUUID),
				isFavorite: mo.Some(true),
			},
			want: &Filter{
				GroupID:    mo.Some(testGroupUUID),
				IsFavorite: mo.Some(true),
			},
			wantErr: false,
		},
		{
			name: "only group and status",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
				status:  mo.Some(StatusPending),
			},
			want: &Filter{
				GroupID: mo.Some(testGroupUUID),
				Status:  mo.Some(StatusPending),
			},
			wantErr: false,
		},
		{
			name: "only group and priority",
			fields: fields{
				groupID:  mo.Some(testGroupUUID),
				priority: mo.Some(PriorityLow),
			},
			want: &Filter{
				GroupID:  mo.Some(testGroupUUID),
				Priority: mo.Some(PriorityLow),
			},
			wantErr: false,
		},
		{
			name: "only group and area",
			fields: fields{
				groupID: mo.Some(testGroupUUID),
				area:    mo.Some(testGeoFilter),
			},
			want: &Filter{
				GroupID: mo.Some(testGroupUUID),
				Area:    mo.Some(testGeoFilter),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			got, err := fb.Build()
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

func TestFilterBuilder_SetArea(t *testing.T) {
	testGeoFilter := GeoFilter{
		Point: GeoPoint{
			Latitude:  67.0,
			Longitude: 52.0,
		},
		Radius: 5.0,
	}
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		point  GeoPoint
		radius float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty area",
			fields: fields{},
			args: args{
				point:  testGeoFilter.Point,
				radius: testGeoFilter.Radius,
			},
			want: &FilterBuilder{
				area: mo.Some(testGeoFilter),
			},
		},
		{
			name: "set area",
			fields: fields{
				area: mo.Some(GeoFilter{
					Point: GeoPoint{
						Latitude:  -testGeoFilter.Point.Latitude,
						Longitude: -testGeoFilter.Point.Longitude,
					},
					Radius: -testGeoFilter.Radius,
				}),
			},
			args: args{
				point:  testGeoFilter.Point,
				radius: testGeoFilter.Radius,
			},
			want: &FilterBuilder{
				area: mo.Some(testGeoFilter),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetArea(tt.args.point, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetArea() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterBuilder_SetGroupID(t *testing.T) {
	testGroupUUID := uuid.New()
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty group id",
			fields: fields{},
			args: args{
				id: testGroupUUID,
			},
			want: &FilterBuilder{
				groupID: mo.Some(testGroupUUID),
			},
		},
		{
			name: "set group id",
			fields: fields{
				groupID: mo.Some(uuid.New()),
			},
			args: args{
				id: testGroupUUID,
			},
			want: &FilterBuilder{
				groupID: mo.Some(testGroupUUID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetGroupID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterBuilder_SetIsFavorite(t *testing.T) {
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		isFavorite bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty is favorite",
			fields: fields{},
			args: args{
				isFavorite: false,
			},
			want: &FilterBuilder{
				isFavorite: mo.Some(false),
			},
		},
		{
			name: "set is favorite",
			fields: fields{
				isFavorite: mo.Some(false),
			},
			args: args{
				isFavorite: true,
			},
			want: &FilterBuilder{
				isFavorite: mo.Some(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetIsFavorite(tt.args.isFavorite); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIsFavorite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterBuilder_SetOwnerID(t *testing.T) {
	testOwnerUUID := uuid.New()
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty owner",
			fields: fields{},
			args: args{
				id: testOwnerUUID,
			},
			want: &FilterBuilder{
				ownerID: mo.Some(testOwnerUUID),
			},
		},
		{
			name: "set owner",
			fields: fields{
				ownerID: mo.Some(uuid.New()),
			},
			args: args{
				id: testOwnerUUID,
			},
			want: &FilterBuilder{
				ownerID: mo.Some(testOwnerUUID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetOwnerID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetOwnerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterBuilder_SetPriority(t *testing.T) {
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		priority Priority
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty priority",
			fields: fields{},
			args: args{
				priority: PriorityLow,
			},
			want: &FilterBuilder{
				priority: mo.Some(PriorityLow),
			},
		},
		{
			name: "set priority",
			fields: fields{
				priority: mo.Some(PriorityLow),
			},
			args: args{
				priority: PriorityMedium,
			},
			want: &FilterBuilder{
				priority: mo.Some(PriorityMedium),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetPriority(tt.args.priority); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterBuilder_SetStatus(t *testing.T) {
	type fields struct {
		ownerID    mo.Option[uuid.UUID]
		groupID    mo.Option[uuid.UUID]
		isFavorite mo.Option[bool]
		status     mo.Option[Status]
		priority   mo.Option[Priority]
		area       mo.Option[GeoFilter]
	}
	type args struct {
		status Status
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FilterBuilder
	}{
		{
			name:   "empty status",
			fields: fields{},
			args: args{
				status: StatusPending,
			},
			want: &FilterBuilder{
				status: mo.Some(StatusPending),
			},
		},
		{
			name: "set status",
			fields: fields{
				status: mo.Some(StatusPending),
			},
			args: args{
				status: StatusCompleted,
			},
			want: &FilterBuilder{
				status: mo.Some(StatusCompleted),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FilterBuilder{
				ownerID:    tt.fields.ownerID,
				groupID:    tt.fields.groupID,
				isFavorite: tt.fields.isFavorite,
				status:     tt.fields.status,
				priority:   tt.fields.priority,
				area:       tt.fields.area,
			}
			if got := fb.SetStatus(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Validate(t *testing.T) {
	type fields struct {
		OwnerID    mo.Option[uuid.UUID]
		GroupID    mo.Option[uuid.UUID]
		IsFavorite mo.Option[bool]
		Status     mo.Option[Status]
		Priority   mo.Option[Priority]
		Area       mo.Option[GeoFilter]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty owner",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "set owner",
			fields: fields{
				OwnerID: mo.Some(uuid.New()),
			},
			wantErr: false,
		},
		{
			name:    "empty group",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "set group",
			fields: fields{
				GroupID: mo.Some(uuid.New()),
			},
			wantErr: false,
		},
		{
			name: "owner and group ids",
			fields: fields{
				OwnerID: mo.Some(uuid.New()),
				GroupID: mo.Some(uuid.New()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				OwnerID:    tt.fields.OwnerID,
				GroupID:    tt.fields.GroupID,
				IsFavorite: tt.fields.IsFavorite,
				Status:     tt.fields.Status,
				Priority:   tt.fields.Priority,
				Area:       tt.fields.Area,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

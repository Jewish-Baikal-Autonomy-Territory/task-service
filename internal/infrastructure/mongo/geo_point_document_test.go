package mongo

import (
	"reflect"
	"task-service/internal/domain/task"
	"testing"
)

func Test_fromDomainGeoPoint(t *testing.T) {
	type args struct {
		point task.GeoPoint
	}
	tests := []struct {
		name string
		args args
		want *geoPointDocument
	}{
		{
			name: "valid geo point",
			args: args{
				point: task.GeoPoint{
					Latitude:  67.0,
					Longitude: 52.0,
				},
			},
			want: &geoPointDocument{
				Type: "Point",
				Coordinates: [2]float32{
					52.0,
					67.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDomainGeoPoint(tt.args.point); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainGeoPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_geoPointDocument_toDomain(t *testing.T) {
	type fields struct {
		Type        string
		Coordinates [2]float32
	}
	tests := []struct {
		name   string
		fields fields
		want   task.GeoPoint
	}{
		{
			name: "valid point",
			fields: fields{
				Type: "Point",
				Coordinates: [2]float32{
					67.0,
					52.0,
				},
			},
			want: task.GeoPoint{
				Latitude:  52.0,
				Longitude: 67.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &geoPointDocument{
				Type:        tt.fields.Type,
				Coordinates: tt.fields.Coordinates,
			}
			if got := gd.toDomain(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

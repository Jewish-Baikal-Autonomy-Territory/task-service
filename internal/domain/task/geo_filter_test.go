package task

import (
	"reflect"
	"testing"
)

func TestNewGeoFilter(t *testing.T) {
	type args struct {
		point  GeoPoint
		radius float64
	}
	tests := []struct {
		name    string
		args    args
		want    GeoFilter
		wantErr bool
	}{
		{
			name: "invalid radius",
			args: args{
				point:  GeoPoint{},
				radius: -1.0,
			},
			want:    GeoFilter{},
			wantErr: true,
		},
		{
			name: "valid radius",
			args: args{
				point:  GeoPoint{},
				radius: 10.0,
			},
			want: GeoFilter{
				Point:  GeoPoint{},
				Radius: 10.0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGeoFilter(tt.args.point, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeoFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGeoFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

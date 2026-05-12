package task

import (
	"reflect"
	"testing"
)

func TestNewGeoPoint(t *testing.T) {
	type args struct {
		latitude  float32
		longitude float32
	}
	tests := []struct {
		name    string
		args    args
		want    GeoPoint
		wantErr bool
	}{
		{
			name: "Invalid latitude value",
			args: args{
				latitude:  -91,
				longitude: 0,
			},
			want:    GeoPoint{},
			wantErr: true,
		},
		{
			name: "Invalid longitude value",
			args: args{
				latitude:  0,
				longitude: 181,
			},
			want:    GeoPoint{},
			wantErr: true,
		},
		{
			name: "Valid latitude and longitude values",
			args: args{
				latitude:  50.0,
				longitude: 120.0,
			},
			want: GeoPoint{
				Latitude:  50.0,
				Longitude: 120.0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGeoPoint(tt.args.latitude, tt.args.longitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeoPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGeoPoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}

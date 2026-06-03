package task

import (
	"reflect"
	"testing"
)

func TestNewGeoPoint(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
	}
	tests := []struct {
		name    string
		args    args
		want    GeoPoint
		wantErr bool
	}{
		{
			name: "invalid latitude and longitude",
			args: args{
				latitude:  -100,
				longitude: 1000,
			},
			want:    GeoPoint{},
			wantErr: true,
		},
		{
			name: "invalid latitude",
			args: args{
				latitude:  -100,
				longitude: 20,
			},
			want:    GeoPoint{},
			wantErr: true,
		},
		{
			name: "invalid longitude",
			args: args{
				latitude:  29,
				longitude: 200,
			},
			want:    GeoPoint{},
			wantErr: true,
		},
		{
			name: "valid latitude and longitude",
			args: args{
				latitude:  20,
				longitude: 100,
			},
			want: GeoPoint{
				Latitude:  20,
				Longitude: 100,
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

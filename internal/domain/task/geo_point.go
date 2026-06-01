package task

import "errors"

var (
	ErrInvalidLatitude  = errors.New("invalid latitude")
	ErrInvalidLongitude = errors.New("invalid longitude")
)

type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

func NewGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	if latitude < -90 || latitude > 90 {
		return GeoPoint{}, ErrInvalidLatitude
	}
	if longitude < -180 || longitude > 180 {
		return GeoPoint{}, ErrInvalidLongitude
	}
	return GeoPoint{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

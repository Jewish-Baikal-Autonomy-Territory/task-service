package task

import "errors"

var ErrInvalidGeoFilter = errors.New("invalid geo filter")

type GeoFilter struct {
	Point  GeoPoint
	Radius float64
}

func NewGeoFilter(point GeoPoint, radius float64) (GeoFilter, error) {
	if radius < 0.0 {
		return GeoFilter{}, ErrInvalidGeoFilter
	}
	return GeoFilter{
		Point:  point,
		Radius: radius,
	}, nil
}

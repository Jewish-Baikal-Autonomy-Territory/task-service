package postgres

import (
	geom "github.com/twpayne/go-geom"
)

func newPoint(latitude, longitude float64) *geom.Point {
	return geom.NewPoint(geom.XY).
		MustSetCoords(geom.Coord{longitude, latitude}).
		SetSRID(4326)
}

func fromPoint(point *geom.Point) (float64, float64) {
	if point == nil {
		return 0, 0
	}
	return point.Y(), point.X()
}

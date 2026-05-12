package mongo

import (
	"task-service/internal/domain/task"
)

type geoPointDocument struct {
	Type        string     `bson:"type"`
	Coordinates [2]float32 `bson:"coordinates"`
}

func (gd *geoPointDocument) toDomain() task.GeoPoint {
	return task.GeoPoint{
		Latitude:  gd.Coordinates[1],
		Longitude: gd.Coordinates[0],
	}
}

func fromDomainGeoPoint(point task.GeoPoint) *geoPointDocument {
	return &geoPointDocument{
		Type: "Point",
		Coordinates: [2]float32{
			point.Longitude,
			point.Latitude,
		},
	}
}

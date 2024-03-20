package model

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const EarthRadiusInKm = 6371.0

type Location struct {
	Id 		   	int64
	Latitude   	float64
	Longitude  	float64
}

func (location *Location) BeforeCreate(scope *gorm.DB) error {
	if err := location.Validate(); err != nil {
		return err
	}
	location.Id = int64(uuid.New().ID()) + time.Now().UnixNano()/int64(time.Microsecond)
	return nil
}

func (location *Location) Validate() error {
	if math.Abs(location.Latitude) > 90 {
		return errors.New("invalid latitude")
	}
	if math.Abs(location.Longitude) > 180 {
		return errors.New("invalid longitude")
	}
	return nil
}

func CalculateDistance(pickedLocation, locationOfInterest Location) float64 {
	currentPositionLatInRad := locationOfInterest.Latitude * math.Pi / 180
	currentPointLatInRad := pickedLocation.Latitude * math.Pi / 180

	deltaLatInRad := math.Abs(pickedLocation.Latitude-locationOfInterest.Latitude) * math.Pi / 180
	deltaLongInRad := math.Abs(pickedLocation.Longitude-locationOfInterest.Longitude) * math.Pi / 180

	a := math.Pow(math.Sin(deltaLatInRad/2), 2) +
		math.Cos(currentPointLatInRad)*math.Cos(currentPositionLatInRad)*
			math.Pow(math.Sin(deltaLongInRad/2), 2)

	return 2 * EarthRadiusInKm * math.Asin(math.Sqrt(a))
}
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"math"
)

const EarthRadiusInKm = 6371.0

type Location struct {
	Id        int64   `bson:"_id,omitempty" json:"id"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

func (l *Location) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(l)
}

func (l *Location) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(l)
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

func (l *Location) Scan(value interface{}) error {
	if value == nil {
		return errors.New("Scan: value is nil")
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan: value is not []byte")
	}

	if err := json.Unmarshal(bytes, l); err != nil {
		return err
	}

	return nil
}

func (l Location) Value() (driver.Value, error) {
	if l.Latitude == 0 && l.Longitude == 0 {
		return nil, nil
	}

	data, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	return data, nil
}

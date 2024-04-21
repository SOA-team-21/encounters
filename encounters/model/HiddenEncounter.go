package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type HiddenEncounter struct {
	Encounter
	Image         string	`bson:"image, omitempty" json:"image"`
	PointLocation Location	`bson:"pointLocation" json:"pointLocation"`
}

type HiddenEncounters []*HiddenEncounter

func (he *HiddenEncounter) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(he)
}

func (he *HiddenEncounter) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(he)
}

func (he *HiddenEncounters) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(he)
}

func (hiddenEncounter *HiddenEncounter)  Validate() error {
	if hiddenEncounter.Image == "" {
		return errors.New("invalid image")
	}
	if hiddenEncounter.PointLocation == (Location{}) {
		return errors.New("point location must not be null")
	}
	return nil
}

func (encounter *HiddenEncounter) CanActivateEncounter(username string,  latitude float64, longitude float64) bool {
	for _, participant := range encounter.Participants {
		if participant.Username == username {
			fmt.Println("Participant already activated")
			return false // Already activated
		}
	}
	for _, completer := range encounter.Completers {
		if completer.Username == username {
			fmt.Println("Participant already completed")
			return false // Already completed
		}
	}

	personsLocation := Location{
		Latitude: latitude,
		Longitude: longitude,
	}

	inProximity := CalculateDistance(personsLocation, encounter.Location) * 1000 <= float64(encounter.Radius)
	if inProximity {
		fmt.Println("Encounter can be activated!")
	}

	fmt.Println("Returning inProximity value:", inProximity)
	return inProximity 
}

func (hiddenEncounter *HiddenEncounter) CanSolveEncounter(username string, longitude, latitude float64) bool {
	personsLocation := Location{Longitude: longitude, Latitude: latitude}
	inRange := CalculateDistance(personsLocation, hiddenEncounter.PointLocation) * 1000 <= 25
	fmt.Println("In range: ", inRange)

	if inRange {
		fmt.Println("Encounter can be solved!")
	}

	fmt.Println("Returning inRange value:", inRange)
	return inRange
}

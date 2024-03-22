package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type HiddenEncounter struct {
	Encounter
	Image         string
	PointLocation Location	`gorm:"type:jsonb;"`
}

func (hiddenEncounter *HiddenEncounter) BeforeCreate(scope *gorm.DB) error {
	if err := hiddenEncounter.Validate(); err != nil {
		return err
	}

	return nil
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


func (encounter *HiddenEncounter) Activate(username string, longitude, latitude float64) bool {
	for _, participant := range encounter.Participants {
		if participant.Username == username {
			return false // Already activated
		}
	}
	for _, completer := range encounter.Completers {
		if completer.Username == username {
			return false // Already completed
		}
	}

	personsLocation := Location{
		Longitude: longitude,
		Latitude:  latitude,
	}

	inProximity := CalculateDistance(personsLocation, encounter.Location) * 1000 <= float64(encounter.Radius)
	if inProximity {
		newParticipant := Participant{
			Username: username,
		}
		encounter.Participants = append(encounter.Participants, newParticipant)
	}

	return inProximity 
}

func (hiddenEncounter *HiddenEncounter) Solve(username string, longitude, latitude float64) bool {
	personsLocation := Location{Longitude: longitude, Latitude: latitude}
	inRange := CalculateDistance(personsLocation, hiddenEncounter.PointLocation) * 1000 <= 25
	if inRange {
		now := time.Now()
		completer := Completer{Username: username, CompletionDate: now}
		hiddenEncounter.Completers = append(hiddenEncounter.Completers, completer)

		for i, participant := range hiddenEncounter.Participants {
			println((participant.Username))
			if participant.Username == username {
				hiddenEncounter.Participants = append(hiddenEncounter.Participants[:i], hiddenEncounter.Participants[i+1:]...)
				break
			}
		}

		for _, participant := range hiddenEncounter.Participants {
			println(participant.Username)
		}
	}

	return inRange
}
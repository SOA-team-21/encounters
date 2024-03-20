package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	hiddenEncounter.Id = int64(uuid.New().ID()) + time.Now().UnixNano()/int64(time.Microsecond)
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

func (hiddenEncounter *HiddenEncounter) Solve(username string, longitude, latitude float64) (bool, error) {
	personsLocation := Location{Longitude: longitude, Latitude: latitude}
	inRange := CalculateDistance(personsLocation, hiddenEncounter.PointLocation) * 1000 <= 5
	if inRange {
		now := time.Now()
		completer := Completer{Username: username, CompletionDate: now}
		hiddenEncounter.Completers = append(hiddenEncounter.Completers, completer)

		for i, participant := range hiddenEncounter.Participants {
			if participant.Username == username {
				hiddenEncounter.Participants = append(hiddenEncounter.Participants[:i], hiddenEncounter.Participants[i+1:]...)
				break
			}
		}
	}

	return inRange, nil
}
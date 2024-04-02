package model

import (
	"errors"

	"gorm.io/gorm"
)

type SocialEncounter struct {
	Encounter
	RequiredParticipants   int64	
	CurrentlyInRange       []Participant `gorm:"type:jsonb;"`
}

func (socialEncounter *SocialEncounter) BeforeCreate(scope *gorm.DB) error {
	if err := socialEncounter.Validate(); err != nil {
		return err
	}
	return nil
}

func (socialEncounter *SocialEncounter)  Validate() error {
	if socialEncounter.RequiredParticipants < 1 {
		return errors.New("Exception! Must be above 0")
	}
	if socialEncounter.CurrentlyInRange == nil {
		return errors.New("Exception! Must not be null!")
	}
	return nil
}

func (encounter *SocialEncounter) Activate(username string, longitude, latitude float64) bool {
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

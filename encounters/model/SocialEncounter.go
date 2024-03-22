package model

import (
	"errors"

	"gorm.io/gorm"
)

type SocialEncounter struct {
	Encounter
	RequiredParticipants   int64	
	CurrentlyInRange       []Participant `gorm:"type:jsonb;"`
	SolveResult            []Completer `gorm:"type:jsonb;"`
	ActivateResult         bool
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

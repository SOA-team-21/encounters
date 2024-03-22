package model

import (
	"errors"

	"gorm.io/gorm"
)

type Participant struct {
	Id 		  int64	`json:"id"`
	Username  string
	EncounterId		int64	`json:"encounterId"`
}

func (participant *Participant) BeforeCreate(scope *gorm.DB) error {
	if err := participant.Validate(); err != nil {
		return err
	}
	return nil
}

func (participant *Participant) Validate() error {
	if participant.Username == "" {
		return errors.New("invalid username")
	}
	return nil
}
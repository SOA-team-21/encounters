package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	participant.Id = int64(uuid.New().ID()) + time.Now().UnixNano()/int64(time.Microsecond)
	return nil
}

func (participant *Participant) Validate() error {
	if participant.Username == "" {
		return errors.New("invalid username")
	}
	return nil
}
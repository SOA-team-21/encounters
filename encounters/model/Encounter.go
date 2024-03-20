package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EncounterStatus int

const (
	Draft EncounterStatus = iota
	Active
	Archived
)

type EncounterType int

const (
	Social EncounterStatus = iota
	Hidden
	Misc
)

type Encounter struct {
	Id           int64	`json:"id" gorm:"primaryKey"`
	Name         string
	Description  string
	Location     Location	`gorm:"type:jsonb;"`
	Experience   int64
	Status       EncounterStatus
	Type         EncounterType
	Radius       int64
	Participants []Participant
	Completers   []Completer
}

func (encounter *Encounter) BeforeCreate(scope *gorm.DB) error {
	encounter.Id = int64(uuid.New().ID()) + time.Now().UnixNano()/int64(time.Microsecond)
	return nil
}
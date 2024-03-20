package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Completer struct {
	Id 				int64	`json:"id"`
	Username        string
	CompletionDate  time.Time
	EncounterId		int64	`json:"encounterId"`
}

func (completer *Completer) BeforeCreate(scope *gorm.DB) error {
	if err := completer.Validate(); err != nil {
		return err
	}
	completer.Id = int64(uuid.New().ID()) + time.Now().UnixNano()/int64(time.Microsecond)
	return nil
}

func (completer *Completer) Validate() error {
	if completer.Username == "" {
		return errors.New("invalid username")
	}
	return nil
}
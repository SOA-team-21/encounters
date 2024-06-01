package model

import (
	"encoding/json"
	"errors"
	"io"
)

type Participant struct {
	Id          int64  `bson:"_id,omitempty" json:"id"`
	Username    string `bson:"username" json:"username"`
	EncounterId int64  `bson:"encounterId" json:"encounterId"`
}

func (p *Participant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Participant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (participant *Participant) Validate() error {
	if participant.Username == "" {
		return errors.New("invalid username")
	}
	return nil
}

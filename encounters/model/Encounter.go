package model

import (
	"encoding/json"
	"io"
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
	Id           int64 				`bson:"_id,omitempty" json:"id"`
	Name         string				`bson:"name" json:"name"`
	Description  string				`bson:"description, omitempty" json:"description"`
	Location     Location			`bson:"location,omitempty" json:"location"`
	Experience   int64				`bson:"experience,omitempty" json:"experience"`
	Status       EncounterStatus	`bson:"status,omitempty" json:"status"`
	Type         EncounterType		`bson:"type,omitempty" json:"type"`
	Radius       int64				`bson:"radius,omitempty" json:"radius"`
	Participants []Participant		`bson:"participants,omitempty" json:"participants"`
	Completers   []Completer		`bson:"completers,omitempty" json:"completers"`
}

func (enc *Encounter) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(enc)
}

func (enc *Encounter) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(enc)
}

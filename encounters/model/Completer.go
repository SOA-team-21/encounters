package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Completer struct {
	Id 				primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username        string				`bson:"username" json:"username"`
	CompletionDate  time.Time			`bson:"completionDate" json:"completionDate"`
	EncounterId		int64				`bson:"encounterId" json:"encounterId"`
}

func (c *Completer) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Completer) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}

func (completer *Completer) Validate() error {
	if completer.Username == "" {
		return errors.New("invalid username")
	}
	return nil
}
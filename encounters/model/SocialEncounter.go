package model

import (
	"encoding/json"
	"fmt"
	"io"
)

type SocialEncounter struct {
	Encounter
	RequiredParticipants   int64	`bson:"requiredParticipants" json:"requiredParticipants"`
	CurrentlyInRange       []Participant `bson:"currentlyInRange" json:"currentlyInRange"`
}

type SocialEncounters []*SocialEncounter

func (he *SocialEncounter) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(he)
}

func (he *SocialEncounter) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(he)
}

func (he *SocialEncounters) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(he)
}


func (encounter *SocialEncounter) CanActivateEncounter(username string,  latitude float64, longitude float64) bool {
	for _, participant := range encounter.Participants {
		if participant.Username == username {
			fmt.Println("Participant already activated")
			return false // Already activated
		}
	}
	for _, completer := range encounter.Completers {
		if completer.Username == username {
			fmt.Println("Participant already completed")
			return false // Already completed
		}
	}

	personsLocation := Location{
		Latitude: latitude,
		Longitude: longitude,
	}

	inProximity := CalculateDistance(personsLocation, encounter.Location) * 1000 <= float64(encounter.Radius)
	if inProximity {
		fmt.Println("Encounter can be activated!")
	}

	fmt.Println("Returning inProximity value:", inProximity)
	return inProximity 
}

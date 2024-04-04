package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type SocialEncounterHandler struct {
	SocialEncounterService *service.SocialEncounterService
}

func (handler *SocialEncounterHandler) Get(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Encounter sa id-em %s", id)
	encounter, err := handler.SocialEncounterService.FindSocialEncounter(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(encounter)
}

func (handler *SocialEncounterHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	encounters, err := handler.SocialEncounterService.GetAll()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	for _, encounter := range encounters {
		fmt.Printf("Retrieved encounter: %v\n", encounter.Id)
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(encounters)
}

func (handler *SocialEncounterHandler) Create (writer http.ResponseWriter, req *http.Request) {
	var SocialEncounter model.SocialEncounter
	err := json.NewDecoder(req.Body).Decode(&SocialEncounter)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	socialEncounter, err := handler.SocialEncounterService.Create(&SocialEncounter)
	if err != nil {
		println("Error while creating a new SocialEncoutner")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(socialEncounter)
}

func (handler *SocialEncounterHandler) Activate (writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	encounter, err := handler.SocialEncounterService.Activate(id, ParticipantLocation)
	if err != nil {
		println("Error while updating a Encounter")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(encounter)
}
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

type HiddenEncounterHandler struct {
	HiddenEncounterService *service.HiddenEncounterService
}

func (handler *HiddenEncounterHandler) Get(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Encounter sa id-em %s", id)
	encounter, err := handler.HiddenEncounterService.FindHiddenEncounter(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(encounter)
}

func (handler *HiddenEncounterHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	encounters, err := handler.HiddenEncounterService.GetAll()
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

func (handler *HiddenEncounterHandler) Create (writer http.ResponseWriter, req *http.Request) {
	var HiddenEncounter model.HiddenEncounter
	err := json.NewDecoder(req.Body).Decode(&HiddenEncounter)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	hiddenEncounter, err := handler.HiddenEncounterService.Create(&HiddenEncounter)
	if err != nil {
		println("Error while creating a new HiddenEncoutner")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(hiddenEncounter)
}

func (handler *HiddenEncounterHandler) Activate (writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	encounter, err := handler.HiddenEncounterService.Activate(id, ParticipantLocation)
	if err != nil {
		println("Error while updating a Encounter")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(encounter)
}

func (handler *HiddenEncounterHandler) Solve (writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	encounter, err := handler.HiddenEncounterService.Solve(id, ParticipantLocation)
	if err != nil {
		println("Error while updating a Encounter")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(encounter)
}
package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"net/http"
)

type HiddenEncounterHandler struct {
	HiddenEncounterService *service.HiddenEncounterService
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
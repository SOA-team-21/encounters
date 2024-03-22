package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"net/http"
)

type SocialEncounterHandler struct {
	SocialEncounterService *service.SocialEncounterService
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
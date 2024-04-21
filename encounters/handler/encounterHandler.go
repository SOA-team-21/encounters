package handler

import (
	"context"
	"encoding/json"
	"encounters/model"
	"encounters/repo"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type EncounterHandler struct {
	logger *log.Logger
	repo *repo.EncounterRepo
}

func NewEncounterHandler(l *log.Logger, r *repo.EncounterRepo) *EncounterHandler {
	return &EncounterHandler{l, r}
}

//hiddenEncounter

func (p *EncounterHandler) GetAllHiddenEncounters(rw http.ResponseWriter, h *http.Request) {
	hiddenEncounters, err := p.repo.GetAllHiddenEncounters()
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if hiddenEncounters == nil {
		return
	}

	err = hiddenEncounters.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *EncounterHandler) GetHiddenEncounterById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	hiddenEncounter, err := p.repo.GetHiddenEncounterById(id)
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if hiddenEncounter == nil {
		http.Error(rw, "Hidden encounter with given id not found", http.StatusNotFound)
		p.logger.Printf("Hidden encounter with id: '%s' not found", id)
		return
	}

	err = hiddenEncounter.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *EncounterHandler) PostHiddenEncounter(rw http.ResponseWriter, h *http.Request) {
	hiddenEncounter := h.Context().Value(KeyProduct{}).(*model.HiddenEncounter)
	p.repo.InsertHiddenEncounter(hiddenEncounter)
	rw.WriteHeader(http.StatusCreated)
}

func (handler *EncounterHandler) ActivateHiddenEncoutner(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.repo.ActivateHiddenEncounter(id, ParticipantLocation)
	writer.WriteHeader(http.StatusOK)
}

func (handler *EncounterHandler) SolveHiddenEncoutner(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	handler.repo.SolveHiddenEncounter(id, ParticipantLocation)
	writer.WriteHeader(http.StatusOK)
}


//socialEncounter


func (p *EncounterHandler) GetAllSocialEncounters(rw http.ResponseWriter, h *http.Request) {
	socialEncounters, err := p.repo.GetAllSocialEncounters()
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if socialEncounters == nil {
		return
	}

	err = socialEncounters.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *EncounterHandler) GetSocialEncounterById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	socialEncounter, err := p.repo.GetSocialEncounterById(id)
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if socialEncounter == nil {
		http.Error(rw, "Hidden encounter with given id not found", http.StatusNotFound)
		p.logger.Printf("Hidden encounter with id: '%s' not found", id)
		return
	}

	err = socialEncounter.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *EncounterHandler) PostSocialEncounter(rw http.ResponseWriter, h *http.Request) {
	socialEncounter := h.Context().Value(KeyProduct{}).(*model.SocialEncounter)
	p.repo.InsertSocialEncounter(socialEncounter)
	rw.WriteHeader(http.StatusCreated)
}

func (handler *EncounterHandler) ActivateSocialEncounter(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var ParticipantLocation model.ParticipantLocation
	err := json.NewDecoder(req.Body).Decode(&ParticipantLocation)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.repo.ActivateSocialEncounter(id, ParticipantLocation)
	writer.WriteHeader(http.StatusOK)
}


func (p *EncounterHandler) MiddlewareHiddenEncounterDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		hiddenEncounter := &model.HiddenEncounter{}
		err := hiddenEncounter.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, hiddenEncounter)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}


func (p *EncounterHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (p *EncounterHandler) MiddlewareSocialEncounterDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		socialEncounter := &model.SocialEncounter{}
		err := socialEncounter.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, socialEncounter)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

package service

import (
	"encounters/model"
	"encounters/repo"
	"fmt"
)

type HiddenEncounterService struct {
	Repo *repo.HiddenEncounterRepo
}

func (service *HiddenEncounterService) FindHiddenEncounter (id string) (*model.HiddenEncounter, error) {
	HiddenEncounter, err := service.Repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &HiddenEncounter, nil
}

func (service *HiddenEncounterService) Create (hiddenEncounter *model.HiddenEncounter) (*model.HiddenEncounter, error) {
	err := service.Repo.CreateHiddenEncounter(hiddenEncounter)
	if err != nil {
		return nil, err
	}
	return hiddenEncounter, nil
}

func (service *HiddenEncounterService) GetAll() ([]model.HiddenEncounter, error) {
	points, err := service.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("encounters are not found")
	}
	return points, nil
}

func (service *HiddenEncounterService) Activate(id string, participantLocation model.ParticipantLocation) (*model.HiddenEncounter, error) {
	encounter, err := service.Repo.Get(id)

	if err != nil {
		return nil, err
	}

	 if encounter.Status != model.Active {
		return nil, fmt.Errorf("encounter is not active")
	 }

	result := encounter.Activate(participantLocation.Username, participantLocation.Longitude, participantLocation.Latitude)
	if !result {
		return nil, fmt.Errorf("activation failed") 
	}

	err = service.Repo.UpdateHiddenEncounter(&encounter) 
	if err != nil {
		return nil, err
	}

	 return &encounter, nil
}

func (service *HiddenEncounterService) Solve(encounterID string, location model.ParticipantLocation) (*model.HiddenEncounter, error) {
	hiddenEncounter, err := service.Repo.Get(encounterID)
	if err != nil {
		return nil, err
	}

	success := hiddenEncounter.Solve(location.Username, location.Longitude, location.Latitude)
	if !success {
		println("failed to solve HiddenEncounter")
	}

	err = service.Repo.UpdateHiddenEncounter(&hiddenEncounter)
	if err != nil {
		println("failed to update HiddenEncounter: %v")
	}

	return &hiddenEncounter, nil
}
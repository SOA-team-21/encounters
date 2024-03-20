package service

import (
	"encounters/model"
	"encounters/repo"
)

type HiddenEncounterService struct {
	Repo *repo.HiddenEncounterRepo
}

func (service *HiddenEncounterService) Create (hiddenEncounter *model.HiddenEncounter) (*model.HiddenEncounter, error) {
	err := service.Repo.CreateHiddenEncounter(hiddenEncounter)
	if err != nil {
		return nil, err
	}
	return hiddenEncounter, nil
}
package service

import (
	"encounters/model"
	"encounters/repo"
)

type SocialEncounterService struct {
	Repo *repo.SocialEncounterRepo
}

func (service *SocialEncounterService) Create (socialEncounter *model.SocialEncounter) (*model.SocialEncounter, error) {
	err := service.Repo.CreateSocialEncounter(socialEncounter)
	if err != nil {
		return nil, err
	}
	return socialEncounter, nil
}

func (service *SocialEncounterService) GetAll() ([]model.SocialEncounter, error) {
	encounters, err := service.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("cannot find encounters")
	}
	for i := range encounters {
		populateEncounter(service, &encounters[i])
	}
	return encounters, nil
}
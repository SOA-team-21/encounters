package service

import (
	"encounters/model"
	"encounters/repo"
	"fmt"
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

func (service *SocialEncounterService) FindSocialEncounter (id string) (*model.SocialEncounter, error) {
	SocialEncounter, err := service.Repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &SocialEncounter, nil
}

func (service *SocialEncounterService) GetAll() ([]model.SocialEncounter, error) {
	points, err := service.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("encounters are not found")
	}
	return points, nil
}

func (service *SocialEncounterService) Activate(id string, participantLocation model.ParticipantLocation) (*model.SocialEncounter, error) {
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

	err = service.Repo.UpdateSocialEncounter(&encounter) 
	if err != nil {
		return nil, err
	}

	 return &encounter, nil
}
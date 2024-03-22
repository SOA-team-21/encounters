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
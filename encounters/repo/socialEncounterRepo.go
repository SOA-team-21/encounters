package repo

import (
	"encounters/model"

	"gorm.io/gorm"
)

type SocialEncounterRepo struct {
	DatabaseConnection *gorm.DB
}


func (repo *SocialEncounterRepo) Get (id string) (model.SocialEncounter, error) {
	socialEncounter := model.SocialEncounter{}
	dbResult := repo.DatabaseConnection.First(&socialEncounter, "id = ?", id)
	if dbResult.Error != nil {
		return socialEncounter, dbResult.Error
	}
	return socialEncounter, nil
}

func (repo *SocialEncounterRepo) CreateSocialEncounter(socialEncounter *model.SocialEncounter) error {
    tx := repo.DatabaseConnection.Begin()

    encounter := model.Encounter{
        Name:        socialEncounter.Name,
        Description: socialEncounter.Description,
        Location: socialEncounter.Location,
		Experience: socialEncounter.Experience,
		Status: socialEncounter.Status,
		Type:  socialEncounter.Type,
		Radius: socialEncounter.Radius,
		Participants: socialEncounter.Participants,
		Completers: socialEncounter.Completers,
    }
    if err := tx.Create(&encounter).Error; err != nil {
        tx.Rollback()
        return err
    }

    socialEncounter.Id = encounter.Id
    if err := tx.Create(socialEncounter).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return err
    }

    return nil
}

func (repo *SocialEncounterRepo) GetAll() ([]model.SocialEncounter, error) {
    encounters := []model.SocialEncounter{}
    dbResult := repo.DatabaseConnection.Find(&tours)
    if dbResult.Error != nil {
        return encounters, dbResult.Error
    }
    return encounters, nil
}
package repo

import (
	"encounters/model"

	"gorm.io/gorm"
)

type HiddenEncounterRepo struct {
	DatabaseConnection *gorm.DB
}


func (repo *HiddenEncounterRepo) Get (id string) (model.HiddenEncounter, error) {
	hiddenEncounter := model.HiddenEncounter{}
	dbResult := repo.DatabaseConnection.Preload("Participants").First(&hiddenEncounter, "id = ?", id)
	if dbResult.Error != nil {
		return hiddenEncounter, dbResult.Error
	}
	return hiddenEncounter, nil
}

func (repo *HiddenEncounterRepo) UpdateHiddenEncounter(encounter *model.HiddenEncounter) error {
	dbResult := repo.DatabaseConnection.Save(encounter)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *HiddenEncounterRepo) GetAll() ([]model.HiddenEncounter, error) {
	encounters := []model.HiddenEncounter{}
	dbResult := repo.DatabaseConnection.Preload("Participants").Find(&encounters)
	if dbResult != nil {
		return encounters, dbResult.Error
	}
	return encounters, nil
}


func (repo *HiddenEncounterRepo) CreateHiddenEncounter(hiddenEncounter *model.HiddenEncounter) error {
    tx := repo.DatabaseConnection.Begin()

    encounter := model.Encounter{
        Name:        hiddenEncounter.Name,
        Description: hiddenEncounter.Description,
        Location: hiddenEncounter.Location,
		Experience: hiddenEncounter.Experience,
		Status: hiddenEncounter.Status,
		Type:  hiddenEncounter.Type,
		Radius: hiddenEncounter.Radius,
		Participants: hiddenEncounter.Participants,
		Completers: hiddenEncounter.Completers,
    }
    if err := tx.Create(&encounter).Error; err != nil {
        tx.Rollback()
        return err
    }

    hiddenEncounter.Id = encounter.Id
    if err := tx.Create(hiddenEncounter).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return err
    }

    return nil
}

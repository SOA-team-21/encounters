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
	dbResult := repo.DatabaseConnection.First(&hiddenEncounter, "id = ?", id)
	if dbResult.Error != nil {
		return hiddenEncounter, dbResult.Error
	}
	return hiddenEncounter, nil
}

func (repo *HiddenEncounterRepo) CreateHiddenEncounter (hiddenEncounter *model.HiddenEncounter) error {
	dbResult := repo.DatabaseConnection.Create(hiddenEncounter)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
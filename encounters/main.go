package main

import (
	"encounters/handler"
	"encounters/model"
	"encounters/repo"
	"encounters/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func initDB() *gorm.DB {
	dsn := "user=postgres password=super dbname=soa_encounters host=localhost port=5432 sslmode=disable search_path=encounters"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	err = database.AutoMigrate(&model.Location{}, &model.Encounter{}, &model.HiddenEncounter{}, &model.Completer{}, &model.Participant{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}
	return database
}

func startServer(handler *handler.HiddenEncounterHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/hiddenEncounter/create", handler.Create).Methods("POST")

	println("Server starting")
	log.Fatal((http.ListenAndServe(":8081", router)))
}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	hiddenEncounterRepo := &repo.HiddenEncounterRepo{DatabaseConnection: database}
	hiddenEncounterService := &service.HiddenEncounterService{Repo: hiddenEncounterRepo}
	hiddenEncounterHandler := &handler.HiddenEncounterHandler{HiddenEncounterService: hiddenEncounterService}

	startServer(hiddenEncounterHandler)
}
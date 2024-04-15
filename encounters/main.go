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
	dsn := "user=postgres password=super dbname=soa_encounters host=enc-database port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	err = database.AutoMigrate(&model.Location{}, &model.Encounter{}, &model.HiddenEncounter{}, &model.SocialEncounter{}, &model.Completer{}, &model.Participant{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}
	return database
}

func startServer(handler *handler.HiddenEncounterHandler, socialEncounterHandler *handler.SocialEncounterHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/hiddenEncounter/create", handler.Create).Methods("POST")
	router.HandleFunc("/hiddenEncounter/getHiddenEncounter/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/hiddenEncounter/activate/{id}", handler.Activate).Methods("PUT")
	router.HandleFunc("/hiddenEncounter/solve/{id}", handler.Solve).Methods("PUT")
	router.HandleFunc("/hiddenEncounter", handler.GetAll).Methods("GET")

	//SocialEncounter
	router.HandleFunc("/socialEncounter/create", socialEncounterHandler.Create).Methods("POST")
	router.HandleFunc("/socialEncounter/activate/{id}", socialEncounterHandler.Activate).Methods("PUT")
	router.HandleFunc("/socialEncounter", socialEncounterHandler.GetAll).Methods("GET")

	println("Server starting")
	log.Fatal((http.ListenAndServe(":8082", router)))
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

	socialEncounterRepo := &repo.SocialEncounterRepo{DatabaseConnection: database}
	socialEncounterService := &service.SocialEncounterService{Repo: socialEncounterRepo}
	socialEncounterHandler := &handler.SocialEncounterHandler{SocialEncounterService: socialEncounterService}

	startServer(hiddenEncounterHandler, socialEncounterHandler)
}

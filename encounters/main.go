package main

import (
	"context"
	"encounters/handler"
	"encounters/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[patient-store] ", log.LstdFlags)

	store, err := repo.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	store.Ping()

	encountersHandler := handler.NewEncounterHandler(logger, store)

	router := mux.NewRouter()
	router.Use(encountersHandler.MiddlewareContentTypeSet)

	//hiddenEncounter

	getHiddenRouter := router.Methods(http.MethodGet).Subrouter()
	getHiddenRouter.HandleFunc("/hiddenEncounter", encountersHandler.GetAllHiddenEncounters)

	postHiddenRouter := router.Methods(http.MethodPost).Subrouter()
	postHiddenRouter.HandleFunc("/hiddenEncounter", encountersHandler.PostHiddenEncounter)
	postHiddenRouter.Use(encountersHandler.MiddlewareHiddenEncounterDeserialization)

	getHiddenByIdRouter := router.Methods(http.MethodGet).Subrouter()
	getHiddenByIdRouter.HandleFunc("/hiddenEncounter/{id}", encountersHandler.GetHiddenEncounterById)

	activateHiddenEncounterRouter := router.Methods(http.MethodPatch).Subrouter()
	activateHiddenEncounterRouter.HandleFunc("/hiddenEncounter/activate/{id}", encountersHandler.ActivateHiddenEncoutner)

	solveHiddenEncounterRouter := router.Methods(http.MethodPatch).Subrouter()
	solveHiddenEncounterRouter.HandleFunc("/hiddenEncounter/solve/{id}", encountersHandler.SolveHiddenEncoutner)

	//socialEncounter

	getSocialRouter := router.Methods(http.MethodGet).Subrouter()
	getSocialRouter.HandleFunc("/socialEncounter", encountersHandler.GetAllSocialEncounters)

	postSocialRouter := router.Methods(http.MethodPost).Subrouter()
	postSocialRouter.HandleFunc("/socialEncounter", encountersHandler.PostSocialEncounter)
	postSocialRouter.Use(encountersHandler.MiddlewareSocialEncounterDeserialization)

	getSocialByIdRouter := router.Methods(http.MethodGet).Subrouter()
	getSocialByIdRouter.HandleFunc("/socialEncounter/{id}", encountersHandler.GetSocialEncounterById)

	activateSocialEncounterRouter := router.Methods(http.MethodPatch).Subrouter()
	activateSocialEncounterRouter.HandleFunc("/socialEncounter/activate/{id}", encountersHandler.ActivateSocialEncounter)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
	
}

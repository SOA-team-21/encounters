package main

import (
	"context"
	"encounters/handler"
	"encounters/proto/encounters"
	"encounters/repo"
	"log"
	"os"
	"os/signal"
	"time"

	"fmt"
	"net"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8081"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[patient-store] ", log.LstdFlags)
	ctx := context.Background()

	store, err := repo.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	store.Ping()

	encounterRepo, err := repo.New(ctx, logger)
	encountersHandler := &handler.EncounterHandler{Repo: encounterRepo}

	lis, err := net.Listen("tcp", ":87")
	fmt.Println("Running gRPC on port 88")
	if err != nil {
		log.Fatalln(err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(lis)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	fmt.Println("Registered gRPC server")

	encounters.RegisterEncountersServiceServer(grpcServer, encountersHandler)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()
	fmt.Println("Serving gRPC")

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}

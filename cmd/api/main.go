package main

import (
	"context"
	"log"

	"github.com/kamil-budzik/csv-processor/internal/api"
	"github.com/kamil-budzik/csv-processor/internal/api/handlers"
	"github.com/kamil-budzik/csv-processor/internal/config"
	"github.com/kamil-budzik/csv-processor/internal/db"
	"github.com/kamil-budzik/csv-processor/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	// DB setup
	database, dbCleanup := db.Setup(cfg)
	defer dbCleanup()

	// Object Storage Setup
	storageClient := storage.Connect(cfg)
	store := storage.NewMinioStorage(storageClient, "csv-uploads")
	ctx := context.Background()
	store.CreateBucket(ctx)

	repo := db.NewTaskRepo(database)
	handler := handlers.NewHandler(repo)
	server := api.NewServer(cfg.Port, handler)

	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

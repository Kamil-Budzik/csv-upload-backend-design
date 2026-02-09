package main

import (
	"log"

	"github.com/kamil-budzik/csv-processor/internal/api"
	"github.com/kamil-budzik/csv-processor/internal/api/handlers"
	"github.com/kamil-budzik/csv-processor/internal/config"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	// DB setup
	database, dbCleanup := db.Setup(cfg)
	defer dbCleanup()

	repo := db.NewTaskRepo(database)
	handler := handlers.NewHandler(repo)
	server := api.NewServer(cfg.Port, handler)

	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

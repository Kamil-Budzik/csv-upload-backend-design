package main

import (
	"log"

	"github.com/kamil-budzik/csv-processor/internal/api"
	"github.com/kamil-budzik/csv-processor/internal/config"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	// DB setup
	dbCleanup := db.Setup(cfg)
	defer dbCleanup()

	server := api.NewServer(cfg.Port)
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

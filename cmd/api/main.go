package main

import (
	"log"

	"github.com/kamil-budzik/csv-processor/internal/api"
	"github.com/kamil-budzik/csv-processor/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	server := api.NewServer(cfg.Port)

	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

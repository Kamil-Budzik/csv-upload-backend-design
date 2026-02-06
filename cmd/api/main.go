package main

import (
	"log"
	"time"

	"github.com/kamil-budzik/csv-processor/internal/api"
	"github.com/kamil-budzik/csv-processor/internal/config"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	// DB Setup
	db.Connect(cfg)
	defer db.DB.Close()
	db.InitDB()
	db.DB.SetMaxOpenConns(25)
	db.DB.SetMaxIdleConns(25)
	db.DB.SetConnMaxLifetime(time.Hour)

	server := api.NewServer(cfg.Port)

	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/kamil-budzik/csv-processor/internal/config"
	_ "github.com/lib/pq" // Postgres driver
)

var DB *sql.DB

func Setup(cfg config.Config) func() {
	Connect(cfg)
	InitDB()
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(time.Hour)

	return func() { DB.Close() }
}

func InitDB() {
	InitTasksTable()
}

func Connect(cfg config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")
}

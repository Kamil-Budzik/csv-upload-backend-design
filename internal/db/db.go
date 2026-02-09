package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/kamil-budzik/csv-processor/internal/config"
	_ "github.com/lib/pq" // Postgres driver
)

func Setup(cfg config.Config) (*sql.DB, func()) {
	conn := Connect(cfg)
	InitDB(conn)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(time.Hour)

	return conn, func() { conn.Close() }
}

func InitDB(DB *sql.DB) {
	InitTasksTable(DB)
}

func Connect(cfg config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return DB
}

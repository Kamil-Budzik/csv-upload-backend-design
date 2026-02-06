package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kamil-budzik/csv-processor/internal/config"
	_ "github.com/lib/pq" // Postgres driver
)

var DB *sql.DB

func InitDB() {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		task_id UUID PRIMARY KEY,
		status VARCHAR(20) CHECK (status IN ('pending', 'processing', 'finished', 'failed')) NOT NULL,
		s3_input_path VARCHAR(255) NOT NULL,
		s3_report_path VARCHAR(255),
		error_message TEXT,
		is_retryable BOOL DEFAULT false NOT NULL,
		created_at TIMESTAMP DEFAULT now() NOT NULL,
		updated_at TIMESTAMP DEFAULT now(),
		original_task_id UUID
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database initialized successfully")
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

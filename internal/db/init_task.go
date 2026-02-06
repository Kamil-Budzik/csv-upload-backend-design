package db

import "log"

func InitTasksTable() {
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

	log.Println("Tasks table initialized successfully")
}

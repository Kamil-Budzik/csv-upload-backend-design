package models

import "time"

type Task struct {
	TaskID         string     `json:"task_id" db:"task_id"`
	Status         string     `json:"status" db:"status"`
	S3InputPath    string     `json:"s3_input_path" db:"s3_input_path"`
	S3ReportPath   *string    `json:"s3_report_path,omitempty" db:"s3_report_path"`
	ErrorMessage   *string    `json:"error_message,omitempty" db:"error_message"`
	IsRetryable    bool       `json:"is_retryable" db:"is_retryable"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	OriginalTaskID *string    `json:"original_task_id,omitempty" db:"original_task_id"`
}

type TaskCreateInput struct {
	S3InputPath string `json:"s3_input_path" binding:"required"`
}

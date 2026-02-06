package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

const taskColumns = `
    task_id,
    status,
    s3_input_path,
    s3_report_path,
    error_message,
    is_retryable,
    created_at,
    updated_at,
    original_task_id
`
const taskTable = "tasks"

// To fetch single row you will need to copy this function but instead of sql.Rows u need to use sql.Row. Annoying but it seems to be true
func scanTask(rows *sql.Rows, t *models.Task) error {
	return rows.Scan(
		&t.TaskID,
		&t.Status,
		&t.S3InputPath,
		&t.S3ReportPath,
		&t.ErrorMessage,
		&t.IsRetryable,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.OriginalTaskID,
	)
}

func GetTasks() ([]models.Task, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s", taskColumns, taskTable)
	rows, err := DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := scanTask(rows, &t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func CreateTask(input models.TaskCreateInput) (models.Task, error) {
	stmt := `
	INSERT INTO tasks (
	    task_id,
	    status,
	    s3_input_path
	)
	VALUES ($1, $2, $3)
	RETURNING task_id, created_at
	`

	var task models.Task
	log.Print("input", input)
	row := DB.QueryRow(stmt, uuid.New(), "pending", input.S3InputPath)
	err := row.Scan(&task.TaskID, &task.CreatedAt)

	if err != nil {
		return task, err
	}

	return task, nil
}

func DeleteTask(id string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE task_id = $1", taskTable)
	_, err := DB.Exec(stmt, id)

	return err

}

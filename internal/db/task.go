package db

import (
	"database/sql"
	"fmt"

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

func scanTasks(rows *sql.Rows, t *models.Task) error {
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

func GetTask(taskId uuid.UUID) (models.Task, error) {
	stmt := fmt.Sprintf("SELECT %s FROM %s WHERE task_id = $1", taskColumns, taskTable)
	row := DB.QueryRow(stmt, taskId)

	var t models.Task

	err := row.Scan(
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

	if err == sql.ErrNoRows {
		return models.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return models.Task{}, err
	}

	return t, nil
}

// Its just a helper for dev env. Later will be protected properly
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
		if err := scanTasks(rows, &t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
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

	row := DB.QueryRow(stmt, uuid.New(), "pending", input.S3InputPath)
	err := row.Scan(&task.TaskID, &task.CreatedAt)

	return task, err
}

func UpdateTask(id uuid.UUID, input models.TaskUpdateStatusInput) (models.Task, error) {
	stmt := `
	UPDATE tasks
	SET status = $1, updated_at = now()
	WHERE task_id = $2
	RETURNING task_id, status, updated_at;
	`

	var task models.Task
	row := DB.QueryRow(stmt, input.Status, id)
	err := row.Scan(&task.TaskID, &task.Status, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.Task{}, ErrTaskNotFound
	}

	return task, err
}

func DeleteTask(id uuid.UUID) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE task_id = $1", taskTable)
	result, err := DB.Exec(stmt, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return err

}

package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	if db == nil {
		panic("NewTaskRepo: db cannot be nil")
	}
	return &TaskRepo{db: db}
}

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

func (r *TaskRepo) GetTask(ctx context.Context, taskId uuid.UUID) (models.Task, error) {
	stmt := fmt.Sprintf("SELECT %s FROM %s WHERE task_id = $1", taskColumns, taskTable)
	row := r.db.QueryRowContext(ctx, stmt, taskId)

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
func (r *TaskRepo) GetTasks(ctx context.Context) ([]models.Task, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s", taskColumns, taskTable)
	rows, err := r.db.QueryContext(ctx, sql)
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

func (r *TaskRepo) CreateTask(ctx context.Context, input models.TaskCreateInput) (models.Task, error) {
	stmt := `
	INSERT INTO tasks (
	    task_id,
	    status,
	    s3_input_path
	)
	VALUES ($1, $2, $3)
	RETURNING task_id, status, s3_input_path, s3_report_path, error_message, is_retryable, created_at, updated_at, original_task_id
	`

	var task models.Task

	row := r.db.QueryRowContext(ctx, stmt, uuid.New(), "pending", input.S3InputPath)
	err := row.Scan(&task.TaskID, &task.Status, &task.S3InputPath, &task.S3ReportPath, &task.ErrorMessage, &task.IsRetryable, &task.CreatedAt, &task.UpdatedAt, &task.OriginalTaskID)

	return task, err
}

func (r *TaskRepo) UpdateTask(ctx context.Context, id uuid.UUID, input models.TaskUpdateStatusInput) (models.Task, error) {
	stmt := `
	UPDATE tasks
	SET status = $1, updated_at = now()
	WHERE task_id = $2
	   AND (
	     (status = 'pending' AND $1 = 'processing')
	     OR (status = 'processing' AND $1 IN ('finished', 'failed'))
	   )
	RETURNING task_id, status, s3_input_path, s3_report_path, error_message, is_retryable, created_at, updated_at, original_task_id
	`

	var task models.Task
	row := r.db.QueryRowContext(ctx, stmt, input.Status, id)
	err := row.Scan(&task.TaskID, &task.Status, &task.S3InputPath, &task.S3ReportPath, &task.ErrorMessage, &task.IsRetryable, &task.CreatedAt, &task.UpdatedAt, &task.OriginalTaskID)

	if err == sql.ErrNoRows {
		var exist string
		checkRow := r.db.QueryRowContext(ctx, "SELECT task_id from tasks WHERE task_id = $1", id)
		new_err := checkRow.Scan(&exist)
		if new_err == sql.ErrNoRows {
			return models.Task{}, ErrTaskNotFound
		}

		return models.Task{}, ErrInvalidTransition
	}

	return task, err
}

func (r *TaskRepo) DeleteTask(ctx context.Context, id uuid.UUID) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE task_id = $1", taskTable)
	result, err := r.db.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return err

}

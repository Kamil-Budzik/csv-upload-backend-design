package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

type TaskRepository interface {
	GetTask(ctx context.Context, id uuid.UUID) (models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	CreateTask(ctx context.Context, input models.TaskCreateInput) (models.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, input models.TaskUpdateStatusInput) (models.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

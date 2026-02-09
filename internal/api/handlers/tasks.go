package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/db"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

type Handler struct {
	repo *db.TaskRepo
}

func NewHandler(repo *db.TaskRepo) *Handler {
	if repo == nil {
		panic("NewHandler: repo cannot be nil")
	}
	return &Handler{repo: repo}
}

func parseUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task_id format",
		})
		return uuid.UUID{}, false
	}
	return id, true
}

func (h *Handler) GetTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	task, err := h.repo.GetTask(id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"task":   task,
	})
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	tasks, err := h.repo.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tasks":  tasks,
	})
}

func (h *Handler) PostTask(c *gin.Context) {
	var input models.TaskCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.repo.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed to Create Task",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Task was created",
		"task":    task,
	})

}

func (h *Handler) PutTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	var input models.TaskUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.repo.UpdateTask(id, input)
	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Failed to update Task",
			"task_id": id,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task was updated",
		"task":    task,
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	err := h.repo.DeleteTask(id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

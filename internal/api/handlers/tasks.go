package handlers

import (
	"errors"
	"fmt"
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
			"status": "error",
			"data":   "Invalid task_id format",
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

	ctx := c.Request.Context()
	task, err := h.repo.GetTask(ctx, id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()
	tasks, err := h.repo.GetTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tasks,
	})
}

func (h *Handler) PostTask(c *gin.Context) {
	var input models.TaskCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": err.Error()})
		return
	}

	ctx := c.Request.Context()
	task, err := h.repo.CreateTask(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   fmt.Sprintf("Failed to Create Task %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   task,
	})

}

func (h *Handler) PutTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	var input models.TaskUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": err.Error()})
		return
	}

	ctx := c.Request.Context()
	task, err := h.repo.UpdateTask(ctx, id, input)
	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   fmt.Sprintf("Failed to update Task %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	ctx := c.Request.Context()
	err := h.repo.DeleteTask(ctx, id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "data": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

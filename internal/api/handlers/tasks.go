package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/db"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

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

func GetTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	task, err := db.GetTask(id)

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

	c.JSON(http.StatusOK, task)
}

func GetAllTasks(c *gin.Context) {
	tasks, err := db.GetTasks()
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

func PostTask(c *gin.Context) {
	var input models.TaskCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := db.CreateTask(input)
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

func PutTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	var input models.TaskUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := db.UpdateTask(id, input)
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

func DeleteTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	err := db.DeleteTask(id)

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

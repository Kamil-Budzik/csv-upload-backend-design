package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/db"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

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
		c.JSON(http.StatusBadRequest, gin.H{
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
	id := c.Param("task_id")
	var input models.TaskUpdateStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := db.UpdateTask(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
	id := c.Param("task_id")

	err := db.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

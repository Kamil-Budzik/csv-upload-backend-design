package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

func GetTasks(c *gin.Context) {
	tasks, err := db.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"taskIds": tasks,
	})
}

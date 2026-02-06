package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

type Task struct {
	TaskID string `json:"task_id"`
}

func GetTasks(c *gin.Context) {
	rows, err := db.DB.Query("SELECT task_id FROM tasks")
	if err != nil {
		c.JSON(500, gin.H{"error": "query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.TaskID); err != nil {
			c.JSON(500, gin.H{"error": "row scan failed: " + err.Error()})
			return
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "rows iteration failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"taskIDs": tasks,
	})
}

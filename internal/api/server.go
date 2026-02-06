package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/db"
)

type Server struct {
	router *gin.Engine
	port   string
}

type Task struct {
	TaskID string `json:"task_id"`
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.router.GET("/", func(c *gin.Context) {
		rows, err := db.DB.Query("select task_id FROM tasks")
		if err != nil {
			c.JSON(500, gin.H{"error": "query failed: " + err.Error()})
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var t Task
			err := rows.Scan(&t.TaskID)
			if err != nil {
				c.JSON(500, gin.H{"error": "row scan failed: " + err.Error()})
			}

			tasks = append(tasks, t)
		}

		if err = rows.Err(); err != nil {
			c.JSON(500, gin.H{"error": "rows iteration failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"taskIDs": tasks,
		})
	})
}

func NewServer(port string) *Server {

	server := &Server{
		router: gin.Default(),
		port:   port,
	}
	server.setupRoutes()

	return server
}

func (s *Server) Run() error {
	return s.router.Run(s.port)
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/api/handlers"
)

type Server struct {
	router *gin.Engine
	port   string
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", handlers.GetHealth)
	s.router.GET("/tasks", handlers.GetTasks)
}

func NewServer(port string) *Server {
	s := &Server{
		router: gin.Default(),
		port:   port,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	return s.router.Run(s.port)
}

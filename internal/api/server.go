package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kamil-budzik/csv-processor/internal/api/handlers"
)

type Server struct {
	router  *gin.Engine
	port    string
	handler *handlers.Handler
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", handlers.GetHealth)
	s.router.GET("/tasks", s.handler.GetAllTasks)
	s.router.GET("/tasks/:task_id", s.handler.GetTask)
	s.router.POST("/tasks", s.handler.PostTask)
	s.router.PUT("/tasks/:task_id", s.handler.PutTask)
	s.router.DELETE("/tasks/:task_id", s.handler.DeleteTask)
}

func NewServer(port string, handler *handlers.Handler) *Server {
	s := &Server{
		router:  gin.Default(),
		port:    port,
		handler: handler,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	return s.router.Run(s.port)
}

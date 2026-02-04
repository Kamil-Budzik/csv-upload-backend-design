package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   string
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

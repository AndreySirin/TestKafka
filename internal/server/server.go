package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(log *slog.Logger, addr string) *Server {
	r := gin.New()
	r.POST("/text", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})
	s := &Server{
		log: log.With("module", "hendler"),
	}
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

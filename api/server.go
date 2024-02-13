package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

// Server serves all HTTP requests for the Tutor Management service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/students", server.createStudent)

	server.router = router

	return server
}

func errorResonse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

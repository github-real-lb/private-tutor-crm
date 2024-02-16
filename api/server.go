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

	// College HTTP Handlers
	router.POST("/colleges", server.createCollege)
	router.GET("/colleges/:id", server.getCollege)
	router.GET("/colleges", server.listColleges)
	router.PUT("/colleges", server.updateCollege)

	// Student HTTP Handlers
	router.POST("/students", server.createStudent)
	router.GET("/students/:id", server.getStudent)
	router.GET("/students", server.listStudents)
	router.PUT("/students", server.updateStudent)
	//router.DELETE("/students/:id", server.deleteStudent) - Transaction

	server.router = router

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResonse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func okResonse(s string) gin.H {
	return gin.H{"message": s}
}

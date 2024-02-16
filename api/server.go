package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

// Server serves all HTTP requests for the Tutor Management service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	// creating the server type with a gin router
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router}

	// adding the colleges HTTP handlers to the router
	router.POST("/colleges", server.createCollege)
	router.GET("/colleges/:id", server.getCollege)
	router.GET("/colleges", server.listColleges)
	router.PUT("/colleges", server.updateCollege)

	// adding the students HTTP handlers to the router
	router.POST("/students", server.createStudent)
	router.GET("/students/:id", server.getStudent)
	router.GET("/students", server.listStudents)
	router.PUT("/students", server.updateStudent)

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

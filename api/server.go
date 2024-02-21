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

	// adding the funnels HTTP handlers to the router
	router.POST("/funnels", server.createFunnel)
	router.GET("/funnels/:id", server.getFunnel)
	router.GET("/funnels", server.listFunnels)
	router.PUT("/funnels", server.updateFunnel)

	// adding the lesson locations HTTP handlers to the router
	router.POST("/lessonlocations", server.createLessonLocation)
	router.GET("/lessonlocations/:id", server.getLessonLocation)
	router.GET("/lessonlocations", server.listLessonLocations)
	router.PUT("/lessonlocations", server.updateLessonLocation)

	// adding the lesson subjects HTTP handlers to the router
	router.POST("/lessonsubjects", server.createLessonSubject)
	router.GET("/lessonsubjects/:id", server.getLessonSubject)
	router.GET("/lessonsubjects", server.listLessonSubjects)
	router.PUT("/lessonsubjects", server.updateLessonSubject)

	// adding the funnels HTTP handlers to the router
	router.POST("/paymentmethods", server.createPaymentMethod)
	router.GET("/paymentmethods/:id", server.getPaymentMethod)
	router.GET("/paymentmethods", server.listPaymentMethods)
	router.PUT("/paymentmethods", server.updatePaymentMethod)

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

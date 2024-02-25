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
	router.POST("/lesson_locations", server.createLessonLocation)
	router.GET("/lesson_locations/:id", server.getLessonLocation)
	router.GET("/lesson_locations", server.listLessonLocations)
	router.PUT("/lesson_locations", server.updateLessonLocation)

	// adding the lesson subjects HTTP handlers to the router
	router.POST("/lesson_subjects", server.createLessonSubject)
	router.GET("/lesson_subjects/:id", server.getLessonSubject)
	router.GET("/lesson_subjects", server.listLessonSubjects)
	router.PUT("/lesson_subjects", server.updateLessonSubject)

	// adding the funnels HTTP handlers to the router
	router.POST("/payment_methods", server.createPaymentMethod)
	router.GET("/payment_methods/:id", server.getPaymentMethod)
	router.GET("/payment_methods", server.listPaymentMethods)
	router.PUT("/payment_methods", server.updatePaymentMethod)

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

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func okResponse(s string) gin.H {
	return gin.H{"message": s}
}

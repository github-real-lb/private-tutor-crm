package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createStudentRequest struct {
	FirstName   string          `json:"first_name" binding:"required"`
	LastName    string          `json:"last_name" binding:"required"`
	Email       sql.NullString  `json:"email"`
	PhoneNumber sql.NullString  `json:"phone_number"`
	Address     sql.NullString  `json:"address"`
	CollegeID   sql.NullInt64   `json:"college_id"`
	FunnelID    sql.NullInt64   `json:"funnel_id"`
	HourlyFee   sql.NullFloat64 `json:"hourly_fee"`
	Notes       sql.NullString  `json:"notes"`
}

func (server *Server) createStudent(ctx *gin.Context) {
	var req createStudentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.CreateStudentParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		CollegeID:   req.CollegeID,
		FunnelID:    req.FunnelID,
		HourlyFee:   req.HourlyFee,
		Notes:       req.Notes,
	}

	student, err := server.store.CreateStudent(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, student)
}

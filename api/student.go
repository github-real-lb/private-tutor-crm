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

type getStudentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getStudent(ctx *gin.Context) {
	var req getStudentRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	student, err := server.store.GetStudent(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResonse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, student)
}

type listStudentsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listStudents(ctx *gin.Context) {
	var req listStudentsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.ListStudentsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	students, err := server.store.ListStudents(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, students)
}

type updateStudentRequest struct {
	StudentID   int64           `json:"student_id" binding:"required"`
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

func (server *Server) updateStudent(ctx *gin.Context) {
	var req updateStudentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.UpdateStudentParams{
		StudentID:   req.StudentID,
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

	err := server.store.UpdateStudent(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResonse("Student updated successfully"))
}

package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createLessonSubjectRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createLessonSubject(ctx *gin.Context) {
	var req createLessonSubjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.CreateLessonSubject(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type getLessonSubjectRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getLessonSubject(ctx *gin.Context) {
	var req getLessonSubjectRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.GetLessonSubject(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResonse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type listLessonSubjectsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listLessonSubjects(ctx *gin.Context) {
	var req listLessonSubjectsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.ListLessonSubjectsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	colleges, err := server.store.ListLessonSubjects(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, colleges)
}

type updateLessonSubjectRequest struct {
	LessonSubjectID int64  `json:"student_id" binding:"required"`
	Name            string `json:"name" binding:"required"`
}

func (server *Server) updateLessonSubject(ctx *gin.Context) {
	var req updateLessonSubjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.UpdateLessonSubjectParams{
		SubjectID: req.LessonSubjectID,
		Name:      req.Name,
	}

	err := server.store.UpdateLessonSubject(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResonse("LessonSubject updated successfully"))
}

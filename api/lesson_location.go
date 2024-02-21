package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createLessonLocationRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createLessonLocation(ctx *gin.Context) {
	var req createLessonLocationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.CreateLessonLocation(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type getLessonLocationRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getLessonLocation(ctx *gin.Context) {
	var req getLessonLocationRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.GetLessonLocation(ctx, req.ID)

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

type listLessonLocationsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listLessonLocations(ctx *gin.Context) {
	var req listLessonLocationsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.ListLessonLocationsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	colleges, err := server.store.ListLessonLocations(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, colleges)
}

type updateLessonLocationRequest struct {
	LessonLocationID int64  `json:"student_id" binding:"required"`
	Name             string `json:"name" binding:"required"`
}

func (server *Server) updateLessonLocation(ctx *gin.Context) {
	var req updateLessonLocationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.UpdateLessonLocationParams{
		LocationID: req.LessonLocationID,
		Name:       req.Name,
	}

	err := server.store.UpdateLessonLocation(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResonse("LessonLocation updated successfully"))
}

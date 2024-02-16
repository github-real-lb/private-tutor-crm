package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createCollegeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createCollege(ctx *gin.Context) {
	var req createCollegeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.CreateCollege(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type getCollegeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCollege(ctx *gin.Context) {
	var req getCollegeRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.GetCollege(ctx, req.ID)

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

type listCollegesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listColleges(ctx *gin.Context) {
	var req listCollegesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.ListCollegesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	colleges, err := server.store.ListColleges(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, colleges)
}

type updateCollegeRequest struct {
	CollegeID int64  `json:"student_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

func (server *Server) updateCollege(ctx *gin.Context) {
	var req updateCollegeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.UpdateCollegeParams{
		CollegeID: req.CollegeID,
		Name:      req.Name,
	}

	err := server.store.UpdateCollege(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResonse("College updated successfully"))
}

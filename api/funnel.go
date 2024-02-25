package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createFunnelRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createFunnel(ctx *gin.Context) {
	var req createFunnelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	college, err := server.store.CreateFunnel(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type getFunnelRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getFunnel(ctx *gin.Context) {
	var req getFunnelRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	college, err := server.store.GetFunnel(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type listFunnelsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listFunnels(ctx *gin.Context) {
	var req listFunnelsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListFunnelsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	colleges, err := server.store.ListFunnels(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, colleges)
}

type updateFunnelRequest struct {
	FunnelID int64  `json:"funnel_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func (server *Server) updateFunnel(ctx *gin.Context) {
	var req updateFunnelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateFunnelParams{
		FunnelID: req.FunnelID,
		Name:     req.Name,
	}

	err := server.store.UpdateFunnel(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResponse("Funnel updated successfully"))
}

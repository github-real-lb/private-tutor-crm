package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
)

type createPaymentMethodRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createPaymentMethod(ctx *gin.Context) {
	var req createPaymentMethodRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.CreatePaymentMethod(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, college)
}

type getPaymentMethodRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPaymentMethod(ctx *gin.Context) {
	var req getPaymentMethodRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	college, err := server.store.GetPaymentMethod(ctx, req.ID)

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

type listPaymentMethodsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPaymentMethods(ctx *gin.Context) {
	var req listPaymentMethodsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.ListPaymentMethodsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	colleges, err := server.store.ListPaymentMethods(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, colleges)
}

type updatePaymentMethodRequest struct {
	PaymentMethodID int64  `json:"student_id" binding:"required"`
	Name            string `json:"name" binding:"required"`
}

func (server *Server) updatePaymentMethod(ctx *gin.Context) {
	var req updatePaymentMethodRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResonse(err))
		return
	}

	arg := db.UpdatePaymentMethodParams{
		PaymentMethodID: req.PaymentMethodID,
		Name:            req.Name,
	}

	err := server.store.UpdatePaymentMethod(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResonse(err))
		return
	}

	ctx.JSON(http.StatusOK, okResonse("PaymentMethod updated successfully"))
}

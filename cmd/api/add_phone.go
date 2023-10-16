package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
)

type AddPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (server *Server) AddPhone(ctx *gin.Context) {
	var req AddPhoneRequest
	var params UserIdParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddPhoneParams{
		ID:    params.Id,
		Phone: sql.NullString{String: req.Phone, Valid: true},
	}
	err := server.store.AddPhone(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		err = fmt.Errorf("failed to update phone number: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "New phone updated")
	return
}

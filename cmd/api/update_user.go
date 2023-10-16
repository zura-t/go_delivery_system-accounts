package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
)

type UserIdParam struct {
	Id int64 `uri:"id"  binding:"required,min=1"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	var params UserIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:   params.Id,
		Name: req.Name,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		err = fmt.Errorf("failed to update user: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	rsp := ConvertUser(user)

	ctx.JSON(http.StatusOK, rsp)
}

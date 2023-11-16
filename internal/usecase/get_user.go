package usecase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		err = fmt.Errorf("failed to find user: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ConvertUser(user)

	ctx.JSON(http.StatusOK, res)
}

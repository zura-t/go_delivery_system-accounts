package usecase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) DeleteUser(ctx *gin.Context) {
	var req DeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", errorResponse(err))
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		err = fmt.Errorf("failed to delete user: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "user has been deleted")
	return
}

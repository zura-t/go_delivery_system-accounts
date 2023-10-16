package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	User                  *UserResponse `json:"user"`
	AccessToken           string        `json:"access_token"`
	AccessTokenExpiresAt  time.Time     `json:"access_token_expires_at"`
	RefreshToken          string        `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time     `json:"refresh_token_expires_at"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Email)
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

	err = pkg.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, server.config.AccessTokenDuration)
	if err != nil {
		err = fmt.Errorf("failed to create access token: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		err = fmt.Errorf("failed to create refresh token: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := &LoginResponse{
		User:                  ConvertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

package usecase

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type UserResponse struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func ConvertUser(user db.User) *UserResponse {
	return &UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone.String,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExists, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userExists != (db.User{}) {
		err = fmt.Errorf("User with such data already exists.")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		err = fmt.Errorf("failed to create user: %s", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := ConvertUser(user)
	ctx.JSON(http.StatusOK, rsp)
}

package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	"github.com/zura-t/go_delivery_system-accounts/internal/usecase"
	"github.com/zura-t/go_delivery_system-accounts/pkg/logger"
)

type userRoutes struct {
	userUsecase usecase.User
	logger      logger.Interface
}

func (server *Server) newUserRoutes(handler *gin.RouterGroup, userUsecase usecase.User, logger logger.Interface) {
	routes := &userRoutes{userUsecase, logger}

	handler.POST("/users", routes.createUser)
	handler.POST("/login", routes.loginUser)

	handler.GET("/users/my_profile", routes.getMyProfile)
	handler.PATCH("/users/", routes.updateUser)
	handler.PATCH("/users/phone_number/", routes.addPhone)
	handler.DELETE("/users/", routes.deleteUser)
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func (r *userRoutes) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.CreateUser(entity.UserRegister{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	User                  entity.User `json:"user"`
}

func (r *userRoutes) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.LoginUser(entity.UserLogin{Email: req.Email, Password: req.Password})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.SetCookie("refresh_token", user.RefreshToken, int(time.Until(user.RefreshTokenExpiresAt).Seconds()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, user)
}

type GetUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (r *userRoutes) getMyProfile(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.GetMyProfile(req.Id)
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type UserIdParam struct {
	Id int64 `uri:"id"  binding:"required,min=1"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (r *userRoutes) updateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var params UserIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, st, err := r.userUsecase.UpdateUser(params.Id, entity.UserUpdate{
		Name: req.Name,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type AddPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (r *userRoutes) addPhone(ctx *gin.Context) {
	var req AddPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var params UserIdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp, st, err := r.userUsecase.AddPhone(params.Id, entity.UserAddPhone{
		Phone: req.Phone,
	})
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}


type DeleteUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (r *userRoutes) deleteUser(ctx *gin.Context) {
	var req DeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, st, err := r.userUsecase.DeleteUser(req.Id)
	if err != nil {
		errorResponse(ctx, st, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (r *userRoutes) logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "logged out")
}

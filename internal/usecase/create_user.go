package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func ConvertUser(user db.User) entity.User {
	return entity.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone.String,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}
}

func (uc *UserUseCase) CreateUser(req *entity.UserRegister) (*entity.User, int, error) {
	userExists, err := uc.store.GetUserByEmail(context.Background(), req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, http.StatusInternalServerError, err
	}
	if userExists != (db.User{}) {
		err = fmt.Errorf("User with such data already exists.")
		return nil, http.StatusBadRequest, err
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := uc.store.CreateUser(context.Background(), arg)
	if err != nil {
		err = fmt.Errorf("failed to create user: %s", err)
		return nil, http.StatusInternalServerError, err
	}

	rsp := ConvertUser(user)
	return &rsp, http.StatusOK, nil
}

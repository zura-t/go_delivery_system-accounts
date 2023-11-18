package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func (uc *UserUseCase) UpdateUser(id int64, req *entity.UserUpdate) (*entity.User, int, error) {
	arg := db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
	}

	user, err := uc.store.UpdateUser(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return nil, http.StatusNotFound, err
		}
		err = fmt.Errorf("failed to update user: %s", err)
		return nil, http.StatusInternalServerError, err
	}
	rsp := ConvertUser(user)

	return &rsp, http.StatusOK, nil
}

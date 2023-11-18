package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
)

func (uc *UserUseCase) GetUser(id int64) (*entity.User, int, error) {
	user, err := uc.store.GetUser(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return nil, http.StatusNotFound, err
		}

		err = fmt.Errorf("failed to find user: %s", err)
		return nil, http.StatusInternalServerError, err
	}
	res := ConvertUser(user)

	return &res, http.StatusOK, nil
}

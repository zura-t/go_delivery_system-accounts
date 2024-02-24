package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func (uc *UserUseCase) AddAdminRole(id int64) (string, int, error) {
	arg := db.AddAdminRoleParams{
		ID:      id,
		IsAdmin: true,
	}

	err := uc.store.AddAdminRole(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return "", http.StatusNotFound, err
		}
		err = fmt.Errorf("failed to update user: %s", err)
		return "", http.StatusInternalServerError, err
	}

	return "User data has been updated", http.StatusOK, nil
}

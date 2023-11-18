package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func (uc *UserUseCase) AddPhone(id int64, req *entity.UserAddPhone) (string, int, error) {
	arg := db.AddPhoneParams{
		ID:    id,
		Phone: sql.NullString{String: req.Phone, Valid: true},
	}
	err := uc.store.AddPhone(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return "", http.StatusNotFound, err
		}

		err = fmt.Errorf("failed to update phone number: %s", err)
		return "", http.StatusInternalServerError, err
	}

	return "New phone updated", http.StatusOK, nil
}

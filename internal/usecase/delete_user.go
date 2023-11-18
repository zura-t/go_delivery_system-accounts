package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
)

func (uc *UserUseCase) DeleteUser(id int64) (string, int, error) {
	err := uc.store.DeleteUser(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return "", http.StatusNotFound, err
		}
		err = fmt.Errorf("failed to delete user: %s", err)
		return "", http.StatusInternalServerError, err
	}

	return "user has been deleted", http.StatusOK, nil
}

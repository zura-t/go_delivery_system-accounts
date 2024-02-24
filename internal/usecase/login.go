package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
)

func (uc *UserUseCase) LoginUser(req *entity.UserLogin) (*entity.UserLoginResponse, int, error) {
	user, err := uc.store.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %s", err)
			return nil, http.StatusNotFound, err
		}
		err = fmt.Errorf("failed to find user: %s", err)
		return nil, http.StatusInternalServerError, err
	}

	err = pkg.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		err = fmt.Errorf("wrong password")
		return nil, http.StatusBadRequest, err
	}

	accessToken, accessPayload, err := uc.tokenMaker.CreateToken(user.ID, user.Email, uc.config.AccessTokenDuration)
	if err != nil {
		err = fmt.Errorf("failed to create access token: %s", err)
		return nil, http.StatusInternalServerError, err
	}

	refreshToken, refreshPayload, err := uc.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		uc.config.RefreshTokenDuration,
	)
	if err != nil {
		err = fmt.Errorf("failed to create refresh token: %s", err)
		return nil, http.StatusInternalServerError, err
	}

	rsp := &entity.UserLoginResponse{
		User:                  ConvertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	return rsp, http.StatusOK, nil
}

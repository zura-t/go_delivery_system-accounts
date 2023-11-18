package usecase

import (
	"github.com/zura-t/go_delivery_system-accounts/config"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/token"
)

type UserUseCase struct {
	store      db.Store
	tokenMaker token.Maker
	config     *config.Config
}

func New(store db.Store, config *config.Config, tokenMaker token.Maker) *UserUseCase {
	return &UserUseCase{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

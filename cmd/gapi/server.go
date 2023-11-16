package gapi

import (
	"fmt"

	"github.com/zura-t/go_delivery_system-accounts/config"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	"github.com/zura-t/go_delivery_system-accounts/token"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	store      db.Store
	config     config.Config
	tokenMaker token.Maker
}

func NewServer(store db.Store, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	return server, nil
}

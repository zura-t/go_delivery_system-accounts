package gapi

import (
	"github.com/zura-t/go_delivery_system-accounts/internal"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pb"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	store  db.Store
	config internal.Config
}

func NewServer(store db.Store, config internal.Config) (*Server, error) {
	server := &Server{store: store, config: config}
	return server, nil
}

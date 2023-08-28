package gapi

import (
	"context"
	"database/sql"

	"github.com/zura-t/go_delivery_system-accounts/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUser(ctx context.Context, req *pb.UserId) (*pb.User, error) {
	user, err := server.store.GetUser(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	res := convertUser(user)

	return res, nil
}
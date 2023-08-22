package gapi

import (
	"context"
	"database/sql"

	"github.com/zura-t/go_delivery_system-accounts/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) DeleteUser(ctx context.Context, req *pb.UserId) (*emptypb.Empty, error) {
	err := server.store.DeleteUser(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %s", err)
	}

	return &emptypb.Empty{}, nil
}
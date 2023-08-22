package gapi

import (
	"context"
	"database/sql"

	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) AddPhone(ctx context.Context, req *pb.AddPhoneRequest) (*emptypb.Empty, error) {
	arg := db.AddPhoneParams{
		ID:    req.GetId(),
		Phone: sql.NullString{String: req.GetPhone(), Valid: true},
	}
	err := server.store.AddPhone(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to update phone number: %s", err)
	}
	return &emptypb.Empty{}, nil
}

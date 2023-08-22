package gapi

import (
	"context"

	"github.com/lib/pq"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
	"github.com/zura-t/go_delivery_system-accounts/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone.String,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := pkg.HashPassword(req.GetPassword())
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		Name:           req.GetName(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, err
			}
		}
		return nil, err
	}

	rsp := convertUser(user)
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateFullName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
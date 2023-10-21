package grpc

import (
	"context"

	pb "github.com/wesleyburlani/go-rest/internal/transport/grpc/pb"
	"github.com/wesleyburlani/go-rest/internal/users"
)

type UserServiceGrpc struct {
	svc *users.Service
}

func NewUserServiceGrpc(svc *users.Service) *UserServiceGrpc {
	return &UserServiceGrpc{svc: svc}
}

func (u *UserServiceGrpc) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := u.svc.Create(ctx, users.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

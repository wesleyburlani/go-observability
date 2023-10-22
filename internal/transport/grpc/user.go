package grpc

import (
	"context"

	pb "github.com/wesleyburlani/go-observability/internal/transport/grpc/pb"
	"github.com/wesleyburlani/go-observability/internal/users"
)

type UserServiceGrpc struct {
	pb.UnimplementedUserServiceServer
	svc *users.Service
}

func NewUserServiceGrpc(svc *users.Service) *UserServiceGrpc {
	return &UserServiceGrpc{svc: svc}
}

func (u *UserServiceGrpc) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := u.svc.Create(ctx, users.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u *UserServiceGrpc) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := u.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

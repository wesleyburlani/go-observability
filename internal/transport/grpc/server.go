package grpc

import (
	"github.com/defval/di"
	pb "github.com/wesleyburlani/go-rest/internal/transport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func CreateGrpcServer(c *di.Container) *grpc.Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	c.Invoke(func(s *UserServiceGrpc) { pb.RegisterUserServiceServer(grpcServer, s) })

	return grpcServer
}

package grpc

import (
	"github.com/defval/di"
	pb "github.com/wesleyburlani/go-observability/internal/transport/grpc/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func CreateGrpcServer(c *di.Container) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	reflection.Register(grpcServer)

	c.Invoke(func(s *UserServiceGrpc) { pb.RegisterUserServiceServer(grpcServer, s) })

	return grpcServer
}

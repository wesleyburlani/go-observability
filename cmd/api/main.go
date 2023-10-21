package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-rest/internal/config"
	"github.com/wesleyburlani/go-rest/internal/di"
	"github.com/wesleyburlani/go-rest/internal/transport/grpc"
	_http "github.com/wesleyburlani/go-rest/internal/transport/http"
	"github.com/wesleyburlani/go-rest/pkg/logger"
	"github.com/wesleyburlani/go-rest/pkg/observability"
	"github.com/wesleyburlani/go-rest/pkg/utils"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDotEnvConfig(".env")
	utils.PanicOnNotNil(err)

	otelShutdown, err := observability.SetupOtel(ctx, observability.OtelConfig{
		ServiceName:              cfg.ServiceName,
		ServiceVersion:           cfg.ServiceVersion,
		OtelExporterOtlpEndpoint: cfg.OtelExporterOtlpEndpoint,
		OtelExporterOtlpInsecure: cfg.OtelExporterOtlpInsecure,
	})
	utils.PanicOnNotNil(err)

	defer func() {
		err = otelShutdown(context.Background())
		utils.PanicOnNotNil(err)
	}()

	container, err := di.BuildContainer(&cfg)
	utils.PanicOnNotNil(err)

	err = container.Invoke(func(c *config.Config, l *logger.Logger) {
		var wg sync.WaitGroup
		wg.Add(1)
		httpAddr := c.HttpAddress
		grpcAddr := ":4000"
		go func() {
			defer wg.Done()
			grpcServer := grpc.CreateGrpcServer(container)
			listener, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				l.With("address", grpcAddr, "error", err).Error(ctx, "error starting grpc server")
			}
			err = grpcServer.Serve(listener)
			if err != nil {
				l.With("address", grpcAddr, "error", err).Error(ctx, "error starting grpc server")
				os.Exit(1)
			}
		}()
		go func() {
			defer wg.Done()
			app := _http.CreateApp(container)
			err = http.ListenAndServe(httpAddr, app)
			if err != nil {
				l.With("address", httpAddr, "error", err).Error(ctx, "error starting http server")
				os.Exit(1)
			}
		}()
		l.With("address", grpcAddr).Info(ctx, "grpc server started")
		l.With("address", httpAddr).Info(ctx, "http server started")
		wg.Wait()
	})
	utils.PanicOnNotNil(err)
}

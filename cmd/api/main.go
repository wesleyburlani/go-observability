package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/internal/di"
	"github.com/wesleyburlani/go-observability/internal/ports/amqp"
	"github.com/wesleyburlani/go-observability/internal/ports/grpc"
	_http "github.com/wesleyburlani/go-observability/internal/ports/http"
	"github.com/wesleyburlani/go-observability/internal/ports/kafka"
	"github.com/wesleyburlani/go-observability/pkg/logger"
	"github.com/wesleyburlani/go-observability/pkg/observability"
	"github.com/wesleyburlani/go-observability/pkg/utils"
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
		httpAddr := c.HttpAddress
		grpcAddr := c.GrpcAddress
		wg.Add(1)
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
		wg.Add(1)
		go func() {
			defer wg.Done()
			app := _http.CreateApp(container)
			err = http.ListenAndServe(httpAddr, app)
			if err != nil {
				l.With("address", httpAddr, "error", err).Error(ctx, "error starting http server")
				os.Exit(1)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := kafka.StartConsume(ctx, container)
			if err != nil {
				l.With("error", err).Error(ctx, "error starting kafka consumer")
				os.Exit(1)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := amqp.StartConsume(ctx, container)
			if err != nil {
				l.With("error", err).Error(ctx, "error starting amqp consumer")
				os.Exit(1)
			}

		}()

		l.With("address", grpcAddr).Info(ctx, "grpc server started")
		l.With("address", httpAddr).Info(ctx, "http server started")
		wg.Wait()
	})
	utils.PanicOnNotNil(err)
}

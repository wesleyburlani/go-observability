package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/internal/di"
	"github.com/wesleyburlani/go-observability/internal/transport/grpc"
	_http "github.com/wesleyburlani/go-observability/internal/transport/http"
	"github.com/wesleyburlani/go-observability/internal/transport/kafka"
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
		wg.Add(1)
		httpAddr := c.HttpAddress
		grpcAddr := c.GrpcAddress
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
		go func() {
			defer wg.Done()
			consumer, err := kafka.CreateConsumer(container)
			if err != nil {
				l.With("error", err).Error(ctx, "error starting kafka consumer")
				os.Exit(1)
			}
			err = consumer.SubscribeTopics([]string{"users"}, nil)
			if err != nil {
				l.With("error", err).Error(ctx, "error subscribing to kafka topic")
				os.Exit(1)
			}
			l.With("topic", "users").Info(ctx, "subscribed to kafka topic")
			l.Info(ctx, "kafka consumer started")
			for {
				msg, err := consumer.ReadMessage(-1)
				if err == nil {
					l.With("message", string(msg.Value)).Info(ctx, "kafka message received")
				} else {
					l.With("error", err).Error(ctx, "error reading kafka message")
				}
			}
		}()
		l.With("address", grpcAddr).Info(ctx, "grpc server started")
		l.With("address", httpAddr).Info(ctx, "http server started")
		wg.Wait()
	})
	utils.PanicOnNotNil(err)
}

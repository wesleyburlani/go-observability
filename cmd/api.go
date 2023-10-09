package main

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-rest/internal/config"
	"github.com/wesleyburlani/go-rest/internal/di"
	_http "github.com/wesleyburlani/go-rest/internal/transport/http"
	"github.com/wesleyburlani/go-rest/pkg/logger"
	"github.com/wesleyburlani/go-rest/pkg/observability"
	"github.com/wesleyburlani/go-rest/pkg/utils"
)

func main() {
	cfg, err := config.LoadDotEnvConfig(".env")
	utils.PanicOnError(err)

	otelShutdown, err := observability.SetupOTelSDK(context.Background(), cfg.ServiceName, cfg.ServiceVersion)
	utils.PanicOnError(err)

	defer func() {
		err = otelShutdown(context.Background())
		utils.PanicOnError(err)
	}()

	container, err := di.BuildContainer(&cfg)
	utils.PanicOnError(err)

	err = container.Invoke(func(c *config.Config, l *logger.Logger) {
		var wg sync.WaitGroup
		wg.Add(1)
		addr := c.HttpAddress
		go func() {
			defer wg.Done()
			app := _http.CreateApp(container)
			err = http.ListenAndServe(addr, app)
			if err != nil {
				l.With("address", addr, "error", err).Error("error starting http server")
				os.Exit(1)
			}
		}()
		l.With("address", addr).Info("server started")
		wg.Wait()
	})
	utils.PanicOnError(err)
}

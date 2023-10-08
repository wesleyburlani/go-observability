package di

import (
	"log/slog"
	"os"

	"github.com/defval/di"

	"github.com/wesleyburlani/go-rest/internal/config"
	pkg_http "github.com/wesleyburlani/go-rest/pkg/http"
	pkg_http_controllers "github.com/wesleyburlani/go-rest/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-rest/pkg/http/middlewares"
)

func BuildContainer() (*di.Container, error) {
	general := di.Options(
		di.Provide(func() *config.Config {
			c, err := config.LoadEnvConfig(".env")
			if err != nil {
				panic(err)
			}
			return &c
		}),
		di.Provide(func() *slog.Logger {
			return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}),
	)

	middlewares := di.Options(
		di.Provide(pkg_http_middlewares.NewLogger, di.As(new(pkg_http.Middleware))),
	)

	controllers := di.Options(
		di.Provide(pkg_http_controllers.NewHealth, di.As(new(pkg_http.Controller))),
	)

	container, err := di.New(general, middlewares, controllers)
	return container, err
}

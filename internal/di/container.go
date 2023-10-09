package di

import (
	"github.com/defval/di"

	"github.com/wesleyburlani/go-rest/internal/config"
	pkg_http "github.com/wesleyburlani/go-rest/pkg/http"
	pkg_http_controllers "github.com/wesleyburlani/go-rest/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-rest/pkg/http/middlewares"
	"github.com/wesleyburlani/go-rest/pkg/logger"
)

func BuildContainer(c *config.Config) (*di.Container, error) {
	general := di.Options(
		di.Provide(func() *config.Config {
			c, err := config.LoadDotEnvConfig(".env")
			if err != nil {
				panic(err)
			}
			return &c
		}),
		di.Provide(logger.NewLogger),
	)

	observability := di.Options(
		di.Provide(pkg_http_controllers.NewHealth, di.As(new(pkg_http.Controller))),
		di.Provide(pkg_http_middlewares.NewLogger, di.As(new(pkg_http.Middleware))),
	)

	container, err := di.New(general, observability)
	return container, err
}

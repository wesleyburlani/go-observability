package di

import (
	"github.com/defval/di"

	"github.com/wesleyburlani/go-rest/internal/config"
	"github.com/wesleyburlani/go-rest/internal/db"
	pkg_http "github.com/wesleyburlani/go-rest/pkg/http"
	pkg_http_controllers "github.com/wesleyburlani/go-rest/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-rest/pkg/http/middlewares"
	"github.com/wesleyburlani/go-rest/pkg/logger"
	"github.com/wesleyburlani/go-rest/pkg/utils"
)

func BuildContainer(c *config.Config) (*di.Container, error) {
	general := di.Options(
		di.Provide(func() *config.Config { return c }),
		di.Provide(func() *logger.Logger {
			level, err := logger.ParseLevel(c.LogLevel)
			utils.PanicOnNotNil(err)
			return logger.NewLogger(logger.Options{Enabled: c.LogEnabled, Level: level})
		}),
	)

	storage := di.Options(
		di.Provide(func() *db.Database {
			d, err := db.NewDatabase(c.DatabaseUrl)
			utils.PanicOnNotNil(err)
			return d
		}),
	)

	observability := di.Options(
		di.Provide(pkg_http_controllers.NewHealth, di.As(new(pkg_http.Controller))),
		di.Provide(pkg_http_middlewares.NewLogger, di.As(new(pkg_http.Middleware))),
	)

	container, err := di.New(general, storage, observability)
	return container, err
}

package di

import (
	"github.com/defval/di"

	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/internal/ports/amqp"
	"github.com/wesleyburlani/go-observability/internal/ports/grpc"
	http_controllers "github.com/wesleyburlani/go-observability/internal/ports/http/controllers"
	"github.com/wesleyburlani/go-observability/internal/ports/kafka/handlers"
	"github.com/wesleyburlani/go-observability/internal/ports/postgres"
	"github.com/wesleyburlani/go-observability/internal/ports/postgres/repositories"
	"github.com/wesleyburlani/go-observability/internal/ports/stdout/observers"
	"github.com/wesleyburlani/go-observability/internal/users"
	pkg_http "github.com/wesleyburlani/go-observability/pkg/http"
	pkg_http_controllers "github.com/wesleyburlani/go-observability/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-observability/pkg/http/middlewares"
	"github.com/wesleyburlani/go-observability/pkg/logger"
	"github.com/wesleyburlani/go-observability/pkg/utils"
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

	connections := di.Options(
		di.Provide(func(logger *logger.Logger) *amqp.ConnectionManager {
			return amqp.NewConnectionManager(c.AmqpUrl, logger)
		}),
	)

	storage := di.Options(
		di.Provide(func() *postgres.Database {
			d, err := postgres.NewDatabase(c.DatabaseUrl)
			utils.PanicOnNotNil(err)
			return d
		}),
	)

	observability := di.Options(
		di.Provide(pkg_http_controllers.NewHealth, di.As(new(pkg_http.Controller))),
		di.Provide(pkg_http_middlewares.NewLogger, di.As(new(pkg_http.Middleware))),
	)

	users := di.Options(
		di.Provide(repositories.NewUserRepository, di.As(new(users.Repository))),
		di.Provide(users.NewService),
		di.Provide(http_controllers.NewUsers, di.As(new(pkg_http.Controller))),
		di.Provide(grpc.NewUserServiceGrpc),
		di.Provide(handlers.NewUserTopicHandler),
		di.Provide(observers.NewUserEventsObserver, di.As(new(users.UserEventsObserver))),
	)

	container, err := di.New(general, connections, storage, observability, users)
	return container, err
}

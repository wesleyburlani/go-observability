package http

import (
	"net/http"

	"github.com/defval/di"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/wesleyburlani/go-rest/internal/config"
	"github.com/wesleyburlani/go-rest/internal/transport/http/controllers"
	pkg_http_controllers "github.com/wesleyburlani/go-rest/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-rest/pkg/http/middlewares"
)

func CreateApp(c *di.Container) http.Handler {
	r := chi.NewRouter()
	c.Invoke(func(c *config.Config) { r.Use(otelchi.Middleware(c.ServiceName, otelchi.WithChiRoutes(r))) })
	c.Invoke(func(m *pkg_http_middlewares.Logger) { r.Use(m.Handle) })
	r.Group(func(r chi.Router) {
		c.Invoke(func(c *pkg_http_controllers.Health) { r.Mount("/health", c.Router()) })
		c.Invoke(func(c *controllers.Users) { r.Mount("/users", c.Router()) })
	})
	return r
}

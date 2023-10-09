package http

import (
	"net/http"

	"github.com/defval/di"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	pkg_http_controllers "github.com/wesleyburlani/go-rest/pkg/http/controllers"
	pkg_http_middlewares "github.com/wesleyburlani/go-rest/pkg/http/middlewares"
)

func CreateApp(c *di.Container) http.Handler {
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("my-server", otelchi.WithChiRoutes(r)))
	c.Invoke(func(m *pkg_http_middlewares.Logger) { r.Use(m.Handle) })
	r.Group(func(r chi.Router) {
		c.Invoke(func(c *pkg_http_controllers.Health) { r.Mount("/health", c.Router()) })
	})
	return r
}

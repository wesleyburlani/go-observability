package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-rest/internal/config"
	"github.com/wesleyburlani/go-rest/internal/di"
	_http "github.com/wesleyburlani/go-rest/internal/transport/http"
)

func main() {
	container, err := di.BuildContainer()
	if err != nil {
		slog.Error("error building container", err)
	}

	err = container.Invoke(func(c *config.Config, l *slog.Logger) {
		var wg sync.WaitGroup
		wg.Add(1)
		addr := c.HttpAddress
		go func() {
			defer wg.Done()
			app := _http.CreateApp(container)
			err = http.ListenAndServe(addr, app)
			if err != nil {
				l.Error("error starting http server", "address", addr, "error", err)
				os.Exit(1)
			}
		}()
		l.Info("server started", "address", addr)
		wg.Wait()
	})
	if err != nil {
		slog.Error("error invoking container", err)
	}
}

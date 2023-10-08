package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/wesleyburlani/go-rest/internal/di"
	_http "github.com/wesleyburlani/go-rest/internal/transport/http"
)

func main() {
	container, err := di.BuildContainer()
	if err != nil {
		slog.Error("error building container", err)
	}

	err = container.Invoke(func(logger *slog.Logger) {
		var wg sync.WaitGroup
		wg.Add(1)
		addr := ":3000"
		go func() {
			defer wg.Done()
			app := _http.CreateApp(*container)
			err = http.ListenAndServe(addr, app)
			if err != nil {
				logger.Error("error starting http server", "address", addr, "error", err)
				os.Exit(1)
			}
		}()
		logger.Info("server started", "address", addr)
		wg.Wait()
	})
	if err != nil {
		slog.Error("error invoking container", err)
	}
}

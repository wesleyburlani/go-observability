package users_test

import (
	"context"
	"testing"

	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/internal/di"
	"github.com/wesleyburlani/go-observability/internal/users"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.LoadDotEnvConfig("../../.env.test")

	if err != nil {
		t.Fatalf("could not load config: %v", err)
	}

	container, err := di.BuildContainer(&cfg)

	if err != nil {
		t.Fatalf("could not build container: %v", err)
	}

	var svc *users.Service
	err = container.Resolve(&svc)

	if err != nil {
		t.Fatalf("could not resolve service: %v", err)
	}

	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			_, err := svc.Get(ctx, 1)

			if err != nil {
				t.Fatalf("could not get user: %v", err)
			}
		})
	})
}

package observers

import (
	"context"

	"github.com/wesleyburlani/go-observability/internal/users"
	"github.com/wesleyburlani/go-observability/pkg/logger"
)

type UserEventsObserver struct {
	logger *logger.Logger
}

func NewUserEventsObserver(l *logger.Logger) *UserEventsObserver {
	return &UserEventsObserver{l}
}

func (u *UserEventsObserver) OnUserCreated(ctx context.Context, user users.User) {
	u.logger.With("user", user).Info(ctx, "user created")
}

func (u *UserEventsObserver) OnUserUpdated(ctx context.Context, user users.User) {
	u.logger.With("user", user).Info(ctx, "user updated")
}

func (u *UserEventsObserver) OnUserDeleted(ctx context.Context, user users.User) {
	u.logger.With("user", user).Info(ctx, "user deleted")
}

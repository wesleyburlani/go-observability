package users

import "context"

type UserEventsObserver interface {
	OnUserCreated(ctx context.Context, user User)
	OnUserUpdated(ctx context.Context, user User)
	OnUserDeleted(ctx context.Context, user User)
}

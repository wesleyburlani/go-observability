package observers

import (
	"context"
	"encoding/json"

	"github.com/wesleyburlani/go-observability/internal/users"
	"github.com/wesleyburlani/go-observability/pkg/amqp"
	"github.com/wesleyburlani/go-observability/pkg/logger"

	amqp_go "github.com/rabbitmq/amqp091-go"
)

type UserEventsObserver struct {
	connManager *amqp.ConnectionManager
	logger      *logger.Logger
}

func NewUserEventsObserver(connManager *amqp.ConnectionManager, l *logger.Logger) *UserEventsObserver {
	return &UserEventsObserver{connManager, l}
}

func (u *UserEventsObserver) OnUserCreated(ctx context.Context, user users.User) {
	json, err := json.Marshal(user)
	if err != nil {
		u.logger.With("error", err).Error(ctx, "error marshalling user")
		return
	}
	u.publishToQueue(ctx, u.getExchange(), json)
}

func (u *UserEventsObserver) OnUserUpdated(ctx context.Context, user users.User) {
}

func (u *UserEventsObserver) OnUserDeleted(ctx context.Context, user users.User) {
}

func (u *UserEventsObserver) publishToQueue(ctx context.Context, exchange string, message []byte) {
	conn := u.connManager.GetConnection(ctx)
	channel, err := conn.Channel()
	if err != nil {
		u.logger.With("error", err).Error(ctx, "error getting channel")
		return
	}

	channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)

	err = channel.PublishWithContext(ctx, exchange, "", false, false, amqp_go.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		u.logger.With("error", err).Error(ctx, "error publishing message")
		return
	}
}

func (u *UserEventsObserver) getExchange() string {
	return "users.events"
}

package amqp

import (
	"context"

	amqp_go "github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *amqp_go.Delivery) error
}

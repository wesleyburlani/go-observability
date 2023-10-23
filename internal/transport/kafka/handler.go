package kafka

import (
	"context"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *ckafka.Message) error
}

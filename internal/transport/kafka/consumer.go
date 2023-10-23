package kafka

import (
	"context"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/defval/di"
	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/internal/transport/kafka/handlers"
	"github.com/wesleyburlani/go-observability/pkg/logger"
)

type KafkaMessage struct {
	Command string `json:"command"`
}

func StartConsume(ctx context.Context, c *di.Container) error {
	err := c.Invoke(func(config *config.Config, l *logger.Logger) {
		consumer, err := ckafka.NewConsumer(&ckafka.ConfigMap{
			"bootstrap.servers": config.KafkaHosts,
			"group.id":          config.ServiceName,
			"auto.offset.reset": "smallest"})
		if err != nil {
			panic(err)
		}

		topics := []string{"users"}
		err = consumer.SubscribeTopics(topics, nil)
		if err != nil {
			l.With("error", err).Error(ctx, "error subscribing to kafka topic")
			os.Exit(1)
		}
		l.With("topics", topics).Info(ctx, "subscribed to kafka topics")
		l.Info(ctx, "kafka consumer started")
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				messageHandler(ctx, msg, c)
			} else {
				l.With("error", err).Error(ctx, "error reading kafka message")
			}
		}
	})
	return err
}

func messageHandler(ctx context.Context, msg *ckafka.Message, c *di.Container) error {
	var l *logger.Logger
	c.Resolve(&l)
	l.With("message", string(msg.Value)).Info(ctx, "kafka message received")
	topic := *msg.TopicPartition.Topic
	switch topic {
	case "users":
		c.Invoke(func(h *handlers.UserTopicHandler) { h.Handle(ctx, msg) })
	default:
		l.With("topic", topic).Error(ctx, "unknown topic")
	}
	return nil
}

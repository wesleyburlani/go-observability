package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/defval/di"
	"github.com/wesleyburlani/go-observability/internal/config"
)

func CreateConsumer(c *di.Container) (*ckafka.Consumer, error) {
	var config *config.Config
	err := c.Resolve(&config)
	if err != nil {
		return nil, err
	}
	consumer, err := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": config.KafkaHosts,
		"group.id":          config.ServiceName,
		"auto.offset.reset": "smallest"})
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

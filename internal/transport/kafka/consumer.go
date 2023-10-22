package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/defval/di"
)

func CreateConsumer(c *di.Container) error {
	consumer, err := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": "host1:9092,host2:9092",
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})
}

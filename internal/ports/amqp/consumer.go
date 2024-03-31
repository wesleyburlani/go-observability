package amqp

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/defval/di"
	"github.com/wesleyburlani/go-observability/internal/config"
	"github.com/wesleyburlani/go-observability/pkg/logger"

	amqp_go "github.com/rabbitmq/amqp091-go"
)

var EXCHANGES = []string{"users"}

func StartConsume(ctx context.Context, c *di.Container) error {
	err := c.Invoke(func(connManager *ConnectionManager, config *config.Config, l *logger.Logger) {
		var wg sync.WaitGroup
		for _, exchange := range EXCHANGES {
			wg.Add(1)
			go func(exchange string) {
				defer wg.Done()
				queue := exchange + "_" + config.ServiceName + "_" + config.ServiceVersion

				var channel *amqp_go.Channel
				queueConsumer := func() {
					for {
						conn := connManager.GetConnection(ctx)
						var err error
						channel, err = conn.Channel()
						if err == nil {
							break
						}
						time.Sleep(3 * time.Second)
					}

					err := channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)

					if err != nil {
						l.With("error", err).With("exchange", exchange).Error(ctx, "error declaring exchange")
						os.Exit(1)
					}

					dlq := queue + "_dlq"
					_, err = channel.QueueDeclare(dlq, true, false, false, false, nil)

					if err != nil {
						l.With("error", err).With("queue", dlq).Error(ctx, "error declaring dlq")
						os.Exit(1)
					}

					l.With("exchange", exchange, "queue", queue).Debug(ctx, "declaring queue to start consuming messages")
					args := amqp_go.Table{}
					args["x-dead-letter-exchange"] = dlq

					q, err := channel.QueueDeclare(queue, true, false, false, false, args)

					if err != nil {
						l.With("error", err).With("queue", queue).Error(ctx, "error declaring queue")
						os.Exit(1)
					}

					err = channel.QueueBind(q.Name, "", exchange, false, nil)

					if err != nil {
						l.With("error", err, "exchange", exchange, "queue", queue).Error(ctx, "error binding queue")
						os.Exit(1)
					}

					msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)

					if err != nil {
						l.With("error", err, "queue", queue).Error(ctx, "error consuming queue")
						os.Exit(1)
					}

					for d := range msgs {
						messageHandler(ctx, exchange, &d, c)
					}
				}

				for {
					queueConsumer()
					<-channel.NotifyClose(make(chan *amqp_go.Error))
				}
			}(exchange)
		}
		wg.Wait()
	})
	return err
}

func messageHandler(ctx context.Context, exchange string, msg *amqp_go.Delivery, c *di.Container) error {
	var l *logger.Logger
	c.Resolve(&l)
	l.With("message", string(msg.Body), "exchange", exchange).Debug(ctx, "received message")

	switch exchange {
	default:
		l.With("exchange", exchange).Warn(ctx, "unknown exchange")
	}
	return nil
}

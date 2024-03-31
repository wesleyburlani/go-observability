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

				var ch *amqp_go.Channel
				consumerManager := func() {
					for {
						conn := connManager.GetConnection(ctx)
						var err error
						ch, err = conn.Channel()
						if err == nil {
							break
						}
						time.Sleep(3 * time.Second)
					}

					err := ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)

					if err != nil {
						l.With("error", err).With("exchange", exchange).Error(ctx, "error declaring exchange")
						os.Exit(1)
					}

					l.With("exchange", exchange, "queue", queue).Info(ctx, "declaring queue to start consuming messages")
					q, err := ch.QueueDeclare(queue, true, false, false, false, nil)

					if err != nil {
						l.With("error", err).With("queue", queue).Error(ctx, "error declaring queue")
						os.Exit(1)
					}

					err = ch.QueueBind(q.Name, "", exchange, false, nil)

					if err != nil {
						l.With("error", err, "exchange", exchange, "queue", queue).Error(ctx, "error binding queue")
						os.Exit(1)
					}

					msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

					if err != nil {
						l.With("error", err, "queue", queue).Error(ctx, "error consuming queue")
						os.Exit(1)
					}

					for d := range msgs {
						l.With("message", string(d.Body), "queue", queue).Info(ctx, "received message")
						messageHandler(ctx, exchange, &d, c)
					}
				}

				consumerManager()

				go func() {
					<-ch.NotifyClose(make(chan *amqp_go.Error))
					consumerManager()
				}()

			}(exchange)
		}
		wg.Wait()
	})
	return err
}

func messageHandler(ctx context.Context, exchange string, msg *amqp_go.Delivery, c *di.Container) error {
	var l *logger.Logger
	c.Resolve(&l)
	l.With("message", string(msg.Body), "exchange", exchange).Info(ctx, "received message")

	switch exchange {
	default:
		l.With("exchange", exchange).Error(ctx, "unknown exchange")
	}
	return nil
}

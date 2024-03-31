package amqp

import (
	"context"
	"fmt"
	"time"

	amqp_go "github.com/rabbitmq/amqp091-go"
	"github.com/wesleyburlani/go-observability/pkg/logger"
)

const delay = 3

type ConnectionManager struct {
	url    string
	logger *logger.Logger
	conn   *amqp_go.Connection
}

func NewConnectionManager(url string, logger *logger.Logger) *ConnectionManager {
	return &ConnectionManager{url: url, logger: logger, conn: nil}
}

func (c *ConnectionManager) GetConnection(ctx context.Context) *amqp_go.Connection {
	if c.conn != nil {
		return c.conn
	}

	c.logger.Debug(ctx, "connecting to amqp")
	for {
		conn, err := amqp_go.Dial(c.url)
		if err == nil {
			c.conn = conn
			break
		}
		c.logger.With("error", err).Error(ctx, fmt.Sprintf("error connecting to amqp, retrying in %d seconds", delay))
		c.conn = nil
		time.Sleep(delay * time.Second)
	}

	go func() {
		err := <-c.conn.NotifyClose(make(chan *amqp_go.Error))
		c.logger.With("error", err).Warn(ctx, "amqp connection closed")
		c.conn = nil
		c.GetConnection(ctx)
	}()

	c.logger.Debug(ctx, "amqp connection established")
	return c.conn
}

package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wesleyburlani/go-observability/internal/users"
	"github.com/wesleyburlani/go-observability/pkg/logger"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type UserTopicHandler struct {
	logger *logger.Logger
	svc    *users.Service
}

type UserTopicMessage struct {
	Command string `json:"command"`
}

type CreateUserMessage struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserTopicHandler(logger *logger.Logger, svc *users.Service) *UserTopicHandler {
	return &UserTopicHandler{
		logger: logger,
		svc:    svc,
	}
}

func (h *UserTopicHandler) Handle(ctx context.Context, msg *ckafka.Message) error {
	h.logger.With("message", string(msg.Value)).Info(ctx, "kafka message received")
	data := UserTopicMessage{}
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		h.logger.With("error", err).Error(ctx, "error unmarshalling kafka message")
		return err
	}

	switch data.Command {
	case "CreateUser":
		err := h.createUser(ctx, msg)
		return err
	default:
		err := errors.New("unknown command")
		h.logger.With("command", data.Command).Error(ctx, "unknown command")
		return err
	}
}

func (h *UserTopicHandler) createUser(ctx context.Context, msg *ckafka.Message) error {
	data := CreateUserMessage{}
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		h.logger.With("error", err).Error(ctx, "error unmarshalling kafka message")
		return err
	}

	_, err = h.svc.Create(ctx, users.User{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		h.logger.With("error", err).Error(ctx, "error creating user")
		return err
	}

	return nil
}

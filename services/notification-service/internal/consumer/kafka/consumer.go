package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"

	"notification-service/config"
	"notification-service/internal/domain/models"
)

var (
	ErrFailedToUnmarshalMessage = "failed to unmarshal message"
	ErrFailedToSendNotification = "failed to send notification"
)

type NotificationService interface {
	SendEmailNotification(n models.EmailNotification) error
}

type Consumer struct {
	logger  *slog.Logger
	reader  *kafka.Reader
	service NotificationService
}

func NewConsumer(cfg *config.Config, service NotificationService, logger *slog.Logger) *Consumer {
	return &Consumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Kafka.Brokers,
			Topic:   cfg.Kafka.Topic,
			GroupID: cfg.Kafka.GroupID,
		}),
		service: service,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var notification models.EmailNotification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			c.logger.Error(ErrFailedToUnmarshalMessage, err)
			continue
		}

		if err := c.service.SendEmailNotification(notification); err != nil {
			c.logger.Error(ErrFailedToSendNotification, err)
			continue
		}
	}
}

func (c *Consumer) Close() error {
	c.logger.Info("closing Kafka consumer")
	return c.reader.Close()
}

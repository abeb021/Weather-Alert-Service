package app

import (
	"log/slog"
	"notification-service/config"
	"notification-service/internal/service"
	smtptransport "notification-service/internal/transport/smtp"
)

type Container struct {
	SMTPServer *smtptransport.Server
}

func NewContainer(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	notificationService := service.NewNotificationService()

	backend := smtptransport.NewBackend(notificationService)

	smtpServer := smtptransport.NewServer(cfg, backend)

	return &Container{
		SMTPServer: smtpServer,
	}, nil
}

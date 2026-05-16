package app

import (
	"log/slog"
	"notification-service/config"
)

type App struct {
	logger    *slog.Logger
	container *Container
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	container, err := NewContainer(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		logger:    logger,
		container: container,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting SMTP server")

	return a.container.SMTPServer.Run()
}

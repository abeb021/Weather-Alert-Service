package cmd

import (
	"errors"
	"notification-service/config"
	"notification-service/internal/app"
	"notification-service/internal/logger"
)

var (
	ErrConfigLoad           = errors.New("failed to load config")
	ErrServerInitialization = errors.New("failed to initialize server")
	ErrServerRun            = errors.New("failed to run server")
)

func main() {
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Info(ErrConfigLoad.Error())
		return
	}

	app, err := app.New(cfg, logger)
	if err != nil {
		logger.Error(ErrServerInitialization.Error(), err)
		return
	}

	if err = app.Run(); err != nil {
		logger.Error(ErrServerRun.Error(), err)
		return
	}
}

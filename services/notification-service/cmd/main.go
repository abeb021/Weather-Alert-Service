package cmd

import (
	"errors"
	"notification-service/config"
	"notification-service/internal/app"
	"notification-service/internal/logger"
)

var (
	ErrConfigLoad = errors.New("failed to load config")
	ErrServerRun  = errors.New("failed to run smtp server")
)

func main() {
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Info(ErrConfigLoad.Error())
		return
	}

	if err = app.Run(cfg, logger); err != nil {
		logger.Info(ErrServerRun.Error())
		return
	}
}

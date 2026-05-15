package main

import (
	"weather-service/config"
	"weather-service/internal/app"
	"weather-service/internal/logger"
	"weather-service/internal/repository/migrations"
	"os"
)

func main() {
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	if err := migrations.Run(cfg.DB.UserURL); err != nil{
		logger.Error("users migrations", "error", err)
		os.Exit(1)
	}

	app.Run(logger, cfg)
}

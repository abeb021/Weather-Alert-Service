package main

import (
	"auth-service/config"
	"auth-service/internal/app"
	"auth-service/internal/logger"
	"auth-service/internal/repository/migrations"
	"os"
)

func main() {
	logger := logger.NewLog()

	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := migrations.Run(cfg.DB.TokenURL)

	app.Run(logger, cfg)
}

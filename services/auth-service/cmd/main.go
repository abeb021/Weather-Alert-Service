package main

import (
	"auth-service/config"
	"auth-service/internal/app"
	"auth-service/internal/logger"
	"auth-service/internal/repository/migrations/tokens"
	"auth-service/internal/repository/migrations/users"
	"os"
)

func main() {
	logger := logger.NewLog()

	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := tokens.Run(cfg.DB.TokenURL); err != nil{
		logger.Logger.Error("token migrations: %v", "error", err)
	}

	if err := users.Run(cfg.DB.UserURL); err != nil{
		logger.Logger.Error("users migrations: %v", "error", err)
	}

	app.Run(logger, cfg)
}

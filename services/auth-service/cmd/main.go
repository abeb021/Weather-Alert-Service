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
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	if err := users.Run(cfg.DB.UserURL); err != nil{
		logger.Error("users migrations", "error", err)
		os.Exit(1)
	}

	if err := tokens.Run(cfg.DB.TokenURL); err != nil{
		logger.Error("token migrations", "error", err)
		os.Exit(1)
	}

	app.Run(logger, cfg)
}

package main

import (
	"weather-service/config"
	"weather-service/internal/app"
	"weather-service/internal/logger"
	"os"
)

func main() {
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	app.Run(logger, cfg)
}

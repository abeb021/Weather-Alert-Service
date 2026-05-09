package main

import (
	"auth-service/internal/app"
	"auth-service/internal/logger"
)

func main() {
	logger := logger.NewLog()

	app.Run(logger)
}

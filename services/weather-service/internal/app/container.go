package app

import (
	"log/slog"
	"weather-service/config"
	"weather-service/internal/api/handlers"
	"weather-service/internal/cache/redis"
	"weather-service/internal/clients/openweather"
)

type Container struct {
	Logger  *slog.Logger
	Handler *handlers.Handler
	Cache   *redis.Cache
	Client  *openweather.Client
}

func NewContainer(logger *slog.Logger, cfg *config.Config) *Container {
	client := openweather.NewClient(cfg.OpenWeatherAPIKey)
	cache := redis.NewCache(cfg.RedisURL)

	handler := handlers.NewHandler(logger, cache, client)

	return &Container{
		Logger:  logger,
		Handler: handler,
	}
}

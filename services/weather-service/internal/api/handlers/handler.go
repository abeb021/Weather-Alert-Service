package handlers

import (
	"log/slog"
	"weather-service/internal/cache"
	"weather-service/internal/clients"
)

type Handler struct {
	log    *slog.Logger
	cache  cache.WeatherCache
	client clients.WeatherFetcher
}

func NewHandler(log *slog.Logger, cache cache.WeatherCache, client clients.WeatherFetcher) *Handler {
	return &Handler{
		log:    log,
		cache:  cache,
		client: client,
	}
}

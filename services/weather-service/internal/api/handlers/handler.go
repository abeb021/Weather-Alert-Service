package handlers

import (
	"context"
	"log/slog"
	"time"
	"weather-service/internal/domain/models"
)

type WeatherCache interface {
	Get(ctx context.Context, city string) (*models.Weather, error)
	Set(ctx context.Context, city string, weather *models.Weather, ttl time.Duration) error
}

type WeatherFetcher interface {
	Fetch(city string) (*models.Weather, error)
}

type Handler struct {
	log    *slog.Logger
	cache  WeatherCache
	client WeatherFetcher
}

func NewHandler(log *slog.Logger, cache WeatherCache, client WeatherFetcher) *Handler {
	return &Handler{
		log:    log,
		cache:  cache,
		client: client,
	}
}

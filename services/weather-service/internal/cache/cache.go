package cache

import (
	"context"
	"time"
	"weather-service/internal/domain/models"
)

type WeatherCache interface {
	Get(ctx context.Context, city string) (*models.Weather, error)
	Set(ctx context.Context, city string, weather *models.Weather, ttl time.Duration) error
}
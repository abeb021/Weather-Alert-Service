package redis

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"weather-service/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache (addr string) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *Cache) Get(ctx context.Context, city string) (*models.Weather, error) {
	key := "weather:" + strings.ToLower(city)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var weather models.Weather
	if err := json.Unmarshal([]byte(data), &weather); err != nil{
		return nil, err
	}

	return &weather, nil
}

func (c *Cache) Set(ctx context.Context, city string, weather *models.Weather, ttl time.Duration) error {
	data, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	key := "weather:" + strings.ToLower(city)
	return c.client.Set(ctx, key, data, ttl).Err()
} 
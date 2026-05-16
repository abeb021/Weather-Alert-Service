package clients

import "weather-service/internal/domain/models"

type WeatherFetcher interface {
	FetchForecast(city string) (*models.Weather, error)
	Fetch(city string) (*models.Weather, error)
}
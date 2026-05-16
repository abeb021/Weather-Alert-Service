package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"weather-service/internal/domain/models"
	domainErrors "weather-service/internal/domain/errors"
)

const baseURL = "https://api.openweathermap.org/data/2.5/weather"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Fetch(city string) (*models.Weather, error) {
	params := url.Values{}
	params.Set("q", city)
	params.Set("appid", c.apiKey)
	params.Set("units", "metric")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domainErrors.ErrAPIUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:
		return nil, domainErrors.ErrCityNotFound
	case http.StatusUnauthorized:
		return nil, domainErrors.ErrInvalidAPIKey
	default:
		return nil, fmt.Errorf("%w: status %d", domainErrors.ErrAPIUnavailable, resp.StatusCode)
	}

	var wd WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&wd); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", domainErrors.ErrAPIUnavailable, err)
	}

	data := &models.Weather{
		City:        wd.Name,
		Temperature: wd.Main.Temp,
		FeelsLike:   wd.Main.FeelsLike,
		Condition:   wd.Weather[0].Main,
		WindSpeed:   wd.Wind.Speed,
		Humidity:    wd.Main.Humidity,
		Timestamp:   time.Now(),
	}

	return data, nil
}

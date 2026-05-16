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

const baseURL = "https://api.openweathermap.org/data/2.5/"

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

type ForecastData struct {
    City struct {
        Name string `json:"name"`
    } `json:"city"`
    List []struct {
        DtTxt string `json:"dt_txt"`
        Main  struct {
            Temp     float64 `json:"temp"`
            Humidity int     `json:"humidity"`
        } `json:"main"`
        Weather []struct {
            Main string `json:"main"`
        } `json:"weather"`
        Wind struct {
            Speed float64 `json:"speed"`
        } `json:"wind"`
        Pop float64 `json:"pop"`
    } `json:"list"`
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

	fullURL := fmt.Sprintf("%s/weather?%s", baseURL, params.Encode())

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domainErrors.ErrAPIUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		//good, skip
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

func (c *Client) FetchForecast(city string) (*models.Weather, error) {
    params := url.Values{}
    params.Set("q", city)
    params.Set("appid", c.apiKey)
    params.Set("units", "metric")

	fullURL := fmt.Sprintf("%s/forecast?%s", baseURL, params.Encode())

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

    var fd ForecastData
    if err := json.NewDecoder(resp.Body).Decode(&fd); err != nil {
        return nil, fmt.Errorf("%w: decode: %v", domainErrors.ErrAPIUnavailable, err)
    }

    // Маппинг ForecastData → models.Weather
    weather := &models.Weather{
        City:      fd.City.Name,
        Timestamp: time.Now(),
    }

    for _, item := range fd.List {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", item.DtTxt)
		if err != nil {
			continue //битая дата, пропускаем
		}
        weather.Forecast = append(weather.Forecast, models.DailyForecast{
            Date:      parsedTime,
            TempDay:   item.Main.Temp,
            Humidity:  item.Main.Humidity,
            Condition: item.Weather[0].Main,
            WindSpeed: item.Wind.Speed,
            RainProb:  item.Pop,
        })
    }

    return weather, nil
}
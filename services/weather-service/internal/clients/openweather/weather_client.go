package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
	domainErrors "weather-service/internal/domain/errors"
	"weather-service/internal/domain/models"
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
	condition := "Unknown"
	if len(wd.Weather) > 0 {
		condition = wd.Weather[0].Main
	}
	data := &models.Weather{
		City:        wd.Name,
		Temperature: wd.Main.Temp,
		FeelsLike:   wd.Main.FeelsLike,
		Condition:   condition,
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

	type slot struct {
		time      time.Time
		temp      float64
		humidity  int
		windSpeed float64
		pop       float64
		condition string
	}

	days := make(map[string][]slot)
	for _, item := range fd.List {
		parsed, err := time.Parse("2006-01-02 15:04:05", item.DtTxt)
		if err != nil {
			continue //битая дата, пропускаем
		}
		dateKey := parsed.Format("2006-01-02")
		cond := "Unknown"
		if len(item.Weather) > 0 {
			cond = item.Weather[0].Main
		}
		days[dateKey] = append(days[dateKey], slot{
			time:      parsed,
			temp:      item.Main.Temp,
			humidity:  item.Main.Humidity,
			windSpeed: item.Wind.Speed,
			pop:       item.Pop,
			condition: cond,
		})
	}

	weather := &models.Weather{
		City:      fd.City.Name,
		Timestamp: time.Now(),
	}

	for dateKey, slots := range days {
		var dayTemp, nightTemp float64
		var dayCount, nightCount int
		var maxWind, maxPop float64
		var sumHumidity int
		condCount := make(map[string]int)

		for _, s := range slots {
			hour := s.time.Hour()
			if hour >= 6 && hour < 18 {
				dayTemp += s.temp
				dayCount++
			} else {
				nightTemp += s.temp
				nightCount++
			}
			if s.windSpeed > maxWind {
				maxWind = s.windSpeed
			}
			if s.pop > maxPop {
				maxPop = s.pop
			}
			sumHumidity += s.humidity
			condCount[s.condition]++
		}

		var avgDay, avgNight float64
		if dayCount > 0 {
			avgDay = dayTemp / float64(dayCount)
		}
		if nightCount > 0 {
			avgNight = nightTemp / float64(nightCount)
		}

		dominantCond := "Unknown"
		maxCount := 0
		for cond, count := range condCount {
			if count > maxCount {
				maxCount = count
				dominantCond = cond
			}
		}

		parsedDate, _ := time.Parse("2006-01-02", dateKey)
		weather.Forecast = append(weather.Forecast, models.DailyForecast{
			Date:      parsedDate,
			TempDay:   avgDay,
			TempNight: avgNight,
			Condition: dominantCond,
			WindSpeed: maxWind,
			Humidity:  sumHumidity / len(slots),
			RainProb:  maxPop,
		})
	}

	sort.Slice(weather.Forecast, func(i, j int) bool {
		return weather.Forecast[i].Date.Before(weather.Forecast[j].Date)
	})
	return weather, nil
}

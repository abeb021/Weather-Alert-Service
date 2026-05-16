package models

import (
	"time"
)

/*
	{
	    "city": "London",
	    "temperature": 15.2,        // текущая (сегодня)
	    "feels_like": 13.8,         // текущая
	    "condition": "Clouds",      // текущая
	    "wind_speed": 4.5,          // текущая
	    "humidity": 72,             // текущая
	    "timestamp": "2026-05-16T10:30:00Z",
	    "forecast": [               // прогноз на 7 дней
	        { "date": "2026-05-17", "temp_day": 18.5, "temp_night": 11.2, ... },
	        { "date": "2026-05-18", "temp_day": 20.1, "temp_night": 13.4, ... },
	        ...
	    ]
	}
*/
type Weather struct {
	City        string          `json:"city"`
	Temperature float64         `json:"temperature"`
	FeelsLike   float64         `json:"feels_like"`
	Condition   string          `json:"condition"`
	WindSpeed   float64         `json:"wind_speed"`
	Humidity    int             `json:"humidity"`
	Timestamp   time.Time       `json:"timestamp"`
	Forecast    []DailyForecast `json:"forecast,omitempty"`
}

type DailyForecast struct {
	Date      time.Time `json:"date"`
	TempDay   float64   `json:"temp_day"`
	TempNight float64   `json:"temp_night"`
	Condition string    `json:"condition"`
	WindSpeed float64   `json:"wind_speed"`
	Humidity  int       `json:"humidity"`
	RainProb  float64   `json:"rain_probability"`
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Weather represents the response for the weather endpoint
type Weather struct {
	Temp        float64  `json:"temp"`
	Condition   string   `json:"condition"`
	Humidity    int      `json:"humidity"`
	WindSpeed   float64  `json:"wind_speed"`
	Code        int      `json:"code"` // WMO weather code
	Forecast3h  Forecast `json:"forecast_3h"`
	ForecastTom Forecast `json:"forecast_tom"`
}

type Forecast struct {
	Temp      float64 `json:"temp"`
	Condition string  `json:"condition"`
	WindSpeed float64 `json:"wind_speed"`
	Precip    float64 `json:"precip"`
}

// OpenMeteoResponse structure for parsing API response
type OpenMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
		IsDay       int     `json:"is_day"`
	} `json:"current_weather"`
	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2m            []float64 `json:"temperature_2m"`
		Relativehumidity2m       []int     `json:"relativehumidity_2m"`
		PrecipitationProbability []int     `json:"precipitation_probability"`
		Weathercode              []int     `json:"weathercode"`
		Windspeed10m             []float64 `json:"windspeed_10m"`
	} `json:"hourly"`
}

// decodeWeatherCode converts WMO code to string condition
func decodeWeatherCode(code int) string {
	switch code {
	case 0:
		return "Clear"
	case 1, 2, 3:
		return "Cloudy"
	case 45, 48:
		return "Fog"
	case 51, 53, 55:
		return "Drizzle"
	case 61, 63, 65:
		return "Rain"
	case 71, 73, 75:
		return "Snow"
	case 80, 81, 82:
		return "Showers"
	case 95, 96, 99:
		return "Storm"
	default:
		return "Unknown"
	}
}

// handleWeather returns real weather data for Aix-les-Bains from Open-Meteo
func handleWeather(w http.ResponseWriter, r *http.Request) {
	// Aix-les-Bains coordinates: 45.6885, 5.9153
	url := "https://api.open-meteo.com/v1/forecast?latitude=45.6885&longitude=5.9153&current_weather=true&hourly=temperature_2m,relativehumidity_2m,precipitation_probability,weathercode,windspeed_10m&timezone=Europe%2FParis&temperature_unit=celsius&windspeed_unit=kmh&precipitation_unit=mm"

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		http.Error(w, "Failed to fetch weather", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var om OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&om); err != nil {
		log.Printf("Error decoding weather: %v", err)
		http.Error(w, "Failed to decode weather", http.StatusInternalServerError)
		return
	}

	// Get current hour index
	now := time.Now()
	currentHour := now.Hour()

	// Prepare response
	weather := Weather{
		Temp:      om.CurrentWeather.Temperature,
		Condition: decodeWeatherCode(om.CurrentWeather.Weathercode),
		Humidity:  0, // Will be filled from hourly
		WindSpeed: om.CurrentWeather.Windspeed,
		Code:      om.CurrentWeather.Weathercode,
	}

	// Find humidity for current hour
	if len(om.Hourly.Relativehumidity2m) > currentHour {
		weather.Humidity = om.Hourly.Relativehumidity2m[currentHour]
	}

	// Forecast +3h
	idx3h := currentHour + 3
	if len(om.Hourly.Temperature2m) > idx3h {
		weather.Forecast3h = Forecast{
			Temp:      om.Hourly.Temperature2m[idx3h],
			Condition: decodeWeatherCode(om.Hourly.Weathercode[idx3h]),
			WindSpeed: om.Hourly.Windspeed10m[idx3h],
			Precip:    float64(om.Hourly.PrecipitationProbability[idx3h]),
		}
	}

	// Forecast Tomorrow (noon)
	idxTom := currentHour + 24 // Same time tomorrow
	if len(om.Hourly.Temperature2m) > idxTom {
		weather.ForecastTom = Forecast{
			Temp:      om.Hourly.Temperature2m[idxTom],
			Condition: decodeWeatherCode(om.Hourly.Weathercode[idxTom]),
			WindSpeed: om.Hourly.Windspeed10m[idxTom],
			Precip:    float64(om.Hourly.PrecipitationProbability[idxTom]),
		}
	}

	json.NewEncoder(w).Encode(weather)
	log.Printf("[%s] GET /api/weather -> %.1f°C, %s (Wind: %.1f km/h) for Aix-les-Bains",
		time.Now().Format("15:04:05"), weather.Temp, weather.Condition, weather.WindSpeed)
}

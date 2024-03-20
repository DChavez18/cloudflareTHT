package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-redis/redis/v8"
)

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"relativehumidity"`
	Pressure    float64 `json:"pressure"`
}

type WeatherController struct {
	WeatherData WeatherData
}

type Period struct {
	Temperature float64 `json:"temperature"`
	Properties  struct {
		RelativeHumidity float64 `json:"relativeHumidity"`
	} `json:"properties"`
	Pressure float64 `json:"pressure"`
}

func (wc *WeatherController) Update(temperature, humidity, pressure float64) {
	wc.WeatherData.Temperature = temperature
	wc.WeatherData.Humidity = humidity
	wc.WeatherData.Pressure = pressure
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	weatherData, err := fetchWeatherData()
	if err != nil {
		panic(err)
	}

	weatherJSON, err := json.Marshal(weatherData)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = rdb.Set(ctx, "weather", weatherJSON, 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Weather data has been cached")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		val, err := rdb.Get(ctx, "weather").Result()
		if err != nil {
			http.Error(w, "Failed to retrieve weather data from Redis cache", http.StatusInternalServerError)
			return
		}

		var weatherDataFromCache WeatherData
		err = json.Unmarshal([]byte(val), &weatherDataFromCache)
		if err != nil {
			http.Error(w, "Failed to parse weather data from Redis cache", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.New("weather").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Weather Data</title>
			</head>
			<body>
				<h1>Weather Data</h1>
				<p>Temperature: {{.Temperature}}</p>
				<p>Humidity: {{.Humidity}}</p>
				<p>Pressure: {{.Pressure}}</p>
			</body>
			</html>
		`))

		err = tmpl.Execute(w, weatherDataFromCache)
		if err != nil {
			http.Error(w, "Failed to render weather data template", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func fetchWeatherData() (WeatherData, error) {

	resp, err := http.Get("http://api.weather.gov/gridpoints/EWX/31,80/forecast/hourly")
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("failed to fetch weather data. Status code: %d", resp.StatusCode)
	}

	var weatherResponse struct {
		Properties struct {
			Periods []Period `json:"periods"`
		} `json:"properties"`
	}

	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to decode JSON weather data: %v", err)
	}

	if len(weatherResponse.Properties.Periods) == 0 {
		return WeatherData{}, fmt.Errorf("no weather data found")
	}

	currentPeriod := weatherResponse.Properties.Periods[0]
	return WeatherData{
		Temperature: currentPeriod.Temperature,
		Humidity:    currentPeriod.Properties.RelativeHumidity,
		Pressure:    0,
	}, nil
}

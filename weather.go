package main

import (
	"fmt"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

type WeatherController struct {
	WeatherData WeatherData
}


func (wc *WeatherController) Update(temperature, humidity, pressure float64) {
	wc.WeatherData.Temperature = temperature
	wc.WeatherData.Humidity = humidity
	wc.WeatherData.Pressure = pressure
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	ctx := context.Background()

	weatherData, err := fetchWeatherData()
	if err != nil {
		panic(err)
	}

	weatherController := WeatherController{}
	weatherController.Update(weatherData.Temperature, weatherData.Humidity, weatherData.Pressure)

	err = rdb.Set(ctx, "weather", fmt.Sprintf("temperature=%.2f humidity=%.2f pressure=%.2f", weatherController.WeatherData.Temperature, weatherController.WeatherData.Humidity, weatherController.WeatherData.Pressure), 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Weather data has been stored in Redis.", fmt.Sprintf("temperature=%.2f humidity=%.2f pressure=%.2f", weatherController.weatherData.temperature, weatherController.weatherData.humidity, weatherController.weatherData.pressure))
}

func fetchWeatherData() (WeatherData, error) {

	resp, err := http.Get("http://api.weather.gov/gridpoints/TOP/31,80/forecast/hourly")
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("failed to fetch weather data. Status code: %d", resp.StatusCode)
	}

	var weatherResponse struct {
		Properties struct {
			Periods []struct {
				Temperature      float64 `json:"temperature"`
				Humidity         float64 `json:"relativeHumidity"`
				WindSpeed        float64 `json:"windSpeed"`
				WindGust         float64 `json:"windGust"`
				WindDirection    string  `json:"windDirection"`
				Icon             string  `json:"icon"`
				ShortForecast    string  `json:"shortForecast"`
				DetailedForecast string  `json:"detailedForecast"`
			} `json:"periods"`
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
		Humidity:    currentPeriod.Humidity,
		Pressure:    0,
	}, nil
}
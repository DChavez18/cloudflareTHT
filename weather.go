package main

import (
	"fmt"
	"context"

	"github.com/go-redis/redis/v8"
)

type WeatherData struct {
	temperature float64
	humidity    float64
	pressure    float64
}

type WeatherController struct {
	weatherData WeatherData
}


func (wc *WeatherController) update(temperature float64, humidity float64, pressure float64) {
	wc.weatherData.temperature = temperature
	wc.weatherData.humidity = humidity
	wc.weatherData.pressure = pressure
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	ctx := context.Background()

	weatherController := WeatherController{}
	weatherController.update(20, 60, 1013)
	fmt.Println(weatherController.weatherData)
}
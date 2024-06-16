package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Survialander/orchestration-service/internal/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

type WeatherService struct{}

type WeatherServiceInterface interface {
	GetWeatherData(ctx context.Context, city string) (WeatherData, error)
}

type WeatherData struct {
	Current struct {
		Temp_c float32
		Temp_f float32
	}
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (s *WeatherService) GetWeatherData(ctx context.Context, city string) (WeatherData, error) {
	tracer := otel.Tracer(viper.GetString("SERVICE"))
	ctx, span := tracer.Start(ctx, "WeatherService.GetWeatherData")
	defer span.End()

	apiKey := viper.GetString("WEATHER_KEY")
	url, _ := url.Parse("https://api.weatherapi.com/v1/current.json")

	qParams := url.Query()
	qParams.Add("q", city)
	qParams.Add("key", apiKey)
	url.RawQuery = qParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return WeatherData{}, err
	}

	client := utils.GetHttpClient()

	response, err := client.Do(req)
	if err != nil {
		return WeatherData{}, err
	}

	var data WeatherData
	err = json.NewDecoder(response.Body).Decode(&data)

	return data, err
}

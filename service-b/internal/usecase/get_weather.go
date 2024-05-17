package usecase

import (
	"context"
	"fmt"

	"github.com/tiagocosta/otel-zipkin-service-b/configs"
	"github.com/tiagocosta/otel-zipkin-service-b/internal/entity"
	"github.com/tiagocosta/otel-zipkin-service-b/internal/pkg/wheatherapi"
	"github.com/tiagocosta/otel-zipkin-service-b/internal/pkg/zipcodeapi"
	"go.opentelemetry.io/otel/trace"
)

var cfg configs.Conf

type GetWeatherInputDTO struct {
	ZipCode string `json:"zipcode"`
}

type GetWeatherOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

type GetWeatherUseCase struct {
	Tracer trace.Tracer
}

func NewGetWeatherUseCase(tracer trace.Tracer) *GetWeatherUseCase {
	cfg = configs.LoadConfig[configs.Conf](".")
	return &GetWeatherUseCase{
		Tracer: tracer,
	}
}

func (uc *GetWeatherUseCase) Execute(ctx context.Context, input GetWeatherInputDTO) (GetWeatherOutputDTO, error) {
	weather, err := entity.NewWeather(input.ZipCode)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	ctx, spanZip := uc.Tracer.Start(ctx, "fiding zip code")
	city, err := zipcodeapi.FindCity(weather.Zip)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	spanZip.End()

	_, spanWeather := uc.Tracer.Start(ctx, "finding weather")
	weather.City = city
	celsius, err := wheatherapi.GetWeather(weather.City, cfg.WeatherAPIKey)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	spanWeather.End()
	weather.FromCelsius(celsius)

	out := GetWeatherOutputDTO{
		City:       weather.City,
		Celsius:    fmt.Sprintf("%.1f", weather.Celsius),
		Fahrenheit: fmt.Sprintf("%.1f", weather.Fahrenheit),
		Kelvin:     fmt.Sprintf("%.1f", weather.Kelvin),
	}

	return out, nil
}

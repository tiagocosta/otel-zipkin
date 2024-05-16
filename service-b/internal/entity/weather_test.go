package entity_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/tiagocosta/otel-zipkin-service-b/internal/entity"
)

func TestNewWeather(t *testing.T) {
	type args struct {
		zip string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Weather
		wantErr bool
	}{
		{
			"Zip code empty. Should return an error",
			args{
				zip: "",
			},
			nil,
			true,
		},
		{
			"Invalid zipcode. Should return an error",
			args{
				zip: "1234567",
			},
			nil,
			true,
		},
		{
			"Valid CEP. Should return a Weather struct",
			args{
				zip: "71218010",
			},
			&entity.Weather{
				Zip: "71218010",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewWeather(tt.args.zip)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWeather() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeather_FromCelsius(t *testing.T) {
	type fields struct {
		Zip  string
		City string
	}
	type args struct {
		celsius float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *entity.Weather
	}{
		{
			"Valid Celsius value. Should result in valid Fahrenheit and Kelvin temperatures",
			fields{
				Zip:  "71218010",
				City: "Brasília",
			},
			args{
				celsius: 25.6,
			},
			&entity.Weather{
				Zip:        "71218010",
				City:       "Brasília",
				Celsius:    25.6,
				Fahrenheit: 78.08,
				Kelvin:     298.60,
			},
		},
		{
			"No Celsius value. Should result in valid Fahrenheit and Kelvin temperatures related to 0 Celsius",
			fields{
				Zip:  "71218010",
				City: "Brasília",
			},
			args{
				celsius: 0.0,
			},
			&entity.Weather{
				Zip:        "71218010",
				City:       "Brasília",
				Celsius:    0.0,
				Fahrenheit: 32.0,
				Kelvin:     273.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &entity.Weather{
				Zip:  tt.fields.Zip,
				City: tt.fields.City,
			}
			w.FromCelsius(tt.args.celsius)

			tolerance := 0.001
			if diff := math.Abs(w.Celsius - tt.want.Celsius); diff > tolerance {
				t.Errorf("expected Celsius = %v, got %v", tt.want.Fahrenheit, w.Fahrenheit)
			}
			if diff := math.Abs(w.Fahrenheit - tt.want.Fahrenheit); diff > tolerance {
				t.Errorf("expected Fahrenheit = %v, got %v", tt.want.Fahrenheit, w.Fahrenheit)
			}
			if diff := math.Abs(w.Kelvin - tt.want.Kelvin); diff > tolerance {
				t.Errorf("expected Kelvin = %v, got %v", tt.want.Kelvin, w.Kelvin)
			}
		})
	}
}

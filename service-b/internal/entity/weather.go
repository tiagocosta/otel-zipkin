package entity

import (
	"fmt"
	"math"
)

type Weather struct {
	Zip        string
	City       string
	Fahrenheit float64
	Celsius    float64
	Kelvin     float64
}

func NewWeather(zip string) (*Weather, error) {
	weather := &Weather{
		Zip: zip,
	}
	err := weather.Validate()
	if err != nil {
		return nil, err
	}
	return weather, nil
}

func (w *Weather) Validate() error {
	if len(w.Zip) != 8 {
		return fmt.Errorf("invalid zipcode")
	}

	return nil
}

func (w *Weather) FromCelsius(celsius float64) {
	w.Celsius = celsius
	w.toFahrenheit()
	w.toKelvin()
}

func (w *Weather) toFahrenheit() {
	x := w.Celsius * 1.8
	w.Fahrenheit = math.Round(x*100)/100 + 32
}

func (w *Weather) toKelvin() {
	w.Kelvin = math.Round(w.Celsius*100)/100 + 273
}

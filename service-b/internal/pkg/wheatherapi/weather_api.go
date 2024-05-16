package wheatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WeatherAPI struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64  `json:"temp_c"`
		Condition struct{} `json:"condition"`
	} `json:"current"`
}

func GetWeather(city string, apiKey string) (float64, error) {
	url := "https://api.weatherapi.com/v1/current.json?q=" + url.PathEscape(city) + "&key=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var w WeatherAPI
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return 0, err
	}

	if (WeatherAPI{}) == w {
		return 0, fmt.Errorf("can not find zipcode")
	}

	return w.Current.TempC, nil
}

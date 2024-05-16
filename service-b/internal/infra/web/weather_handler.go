package web

import (
	"encoding/json"
	"net/http"

	"github.com/tiagocosta/otel-zipkin-service-b/internal/usecase"
)

type WebWeatherHandler struct {
}

func NewWebWeatherHandler() *WebWeatherHandler {
	return &WebWeatherHandler{}
}

func (h *WebWeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	var dto usecase.GetWeatherInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getWeather := usecase.NewGetWeatherUseCase()
	output, err := getWeather.Execute(dto)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "can not find zipcode" {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

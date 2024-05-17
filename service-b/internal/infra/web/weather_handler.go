package web

import (
	"encoding/json"
	"net/http"

	"github.com/tiagocosta/otel-zipkin-service-b/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WebWeatherHandler struct {
	Tracer trace.Tracer
}

func NewWebWeatherHandler(tracer trace.Tracer) *WebWeatherHandler {
	return &WebWeatherHandler{
		Tracer: tracer,
	}
}

func (h *WebWeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	ctx, span := h.Tracer.Start(ctx, "starting service-b handler")
	defer span.End()
	var dto usecase.GetWeatherInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getWeather := usecase.NewGetWeatherUseCase(h.Tracer)
	output, err := getWeather.Execute(ctx, dto)
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

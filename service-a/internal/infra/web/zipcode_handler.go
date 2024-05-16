package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type ZipCodeDTO struct {
	ZipCode string `json:"cep"`
}

type WeatherInputDTO struct {
	ZipCode string `json:"zipcode"`
}

type WeatherOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

type WebZipCodeHandler struct {
	Tracer trace.Tracer
}

func NewWebZipCodeHandler(tracer trace.Tracer) *WebZipCodeHandler {
	return &WebZipCodeHandler{
		Tracer: tracer,
	}
}

func (h *WebZipCodeHandler) ProcessZipCode(w http.ResponseWriter, r *http.Request) {
	// ctx, span := h.Tracer.Start(r.Context(), "process-zip-code", trace.WithSpanKind(trace.SpanKindServer))
	_, span := h.Tracer.Start(r.Context(), "processing zip code")
	// defer span.End()
	var dto ZipCodeDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if len(dto.ZipCode) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	span.End()

	weatherDTO := WeatherInputDTO(dto)
	res, err := h.callServiceB(weatherDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errMsg, _ := io.ReadAll(res.Body)
		http.Error(w, string(errMsg), res.StatusCode)
		return
	}

	var weatherOutputDTO WeatherOutputDTO
	err = json.NewDecoder(res.Body).Decode(&weatherOutputDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(weatherOutputDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebZipCodeHandler) callServiceB(weatherDTO WeatherInputDTO) (*http.Response, error) {
	jsonBody, err := json.Marshal(weatherDTO)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := "http://localhost:8082/weather"
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

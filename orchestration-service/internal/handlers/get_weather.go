package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Survialander/orchestration-service/internal"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherHandler struct {
	CepService     internal.CepServiceInterface
	WeatherService internal.WeatherServiceInterface
}

type WeatherResponse struct {
	City   string  `json:"city"`
	Temp_c float32 `json:"temp_c"`
	Temp_f float32 `json:"temp_f"`
	Temp_k float32 `json:"temp_k"`
}

func NewWeatherHandler(cepService internal.CepServiceInterface, weatherService internal.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		CepService:     cepService,
		WeatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	tracer := otel.Tracer(viper.GetString("SERVICE"))

	ctx, span := tracer.Start(ctx, "GetWeatherHandler")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	cepData, err := h.CepService.GetCepData(ctx, cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if cepData.Localidade == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	weatherData, err := h.WeatherService.GetWeatherData(ctx, cepData.Localidade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := WeatherResponse{
		City:   cepData.Localidade,
		Temp_c: weatherData.Current.Temp_c,
		Temp_f: weatherData.Current.Temp_f,
		Temp_k: weatherData.Current.Temp_c + 273,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherInfoHandler struct{}

type WeatherInfoBody struct {
	Cep interface{}
}

type OrchestrationResponse struct {
	City   string  `json:"city"`
	Temp_c float32 `json:"temp_c"`
	Temp_f float32 `json:"temp_f"`
	Temp_k float32 `json:"temp_k"`
}

func NewWeatherHandler() *WeatherInfoHandler {
	return &WeatherInfoHandler{}
}

func (h *WeatherInfoHandler) GetWeatherInfo(w http.ResponseWriter, r *http.Request) {
	orchestrationUrl := viper.GetString("ORCHESTRATION_SERVICE_URL")
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	tracer := otel.Tracer(viper.GetString("SERVICE"))

	ctx, span := tracer.Start(ctx, "GetWeatherInfoHandler")
	defer span.End()

	var data WeatherInfoBody
	json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	cep, valid := validateCep(data.Cep)

	if !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	client := http.Client{}

	url, _ := url.Parse(orchestrationUrl)
	query := url.Query()
	query.Add("cep", cep)
	url.RawQuery = query.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", url.String(), nil)

	response, err := client.Do(req)

	if response.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(" can not find zipcode"))
		return
	}

	var responseData OrchestrationResponse
	json.NewDecoder(response.Body).Decode(&responseData)
	defer response.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func validateCep(cep interface{}) (string, bool) {
	str, ok := cep.(string)

	if ok && len(str) == 8 {
		return str, ok
	}

	return "", false
}

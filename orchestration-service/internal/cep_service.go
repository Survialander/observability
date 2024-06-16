package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Survialander/orchestration-service/internal/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

type CepService struct{}

type CepServiceInterface interface {
	GetCepData(ctx context.Context, cep string) (CepData, error)
}

type CepData struct {
	Cep         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	Uf          string
}

func NewCepService() *CepService {
	return &CepService{}
}

func (s *CepService) GetCepData(ctx context.Context, cep string) (CepData, error) {
	tracer := otel.Tracer(viper.GetString("SERVICE"))
	ctx, span := tracer.Start(ctx, "CepService.GetCepData")
	defer span.End()

	client := utils.GetHttpClient()

	url := fmt.Sprintf("https://viacep.com.br/ws/%v/json/", cep)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := client.Do(req)

	if err != nil {
		return CepData{}, err
	}

	var data CepData
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return CepData{}, err
	}

	return data, nil
}

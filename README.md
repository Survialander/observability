# Observabilidade com Open Telemetry

## Descrição
Esse sistema recebe um cep como input, e baseado nas informações deste cep retorna a cidade e o clima em tempo real, rastreando as requests com o OpenTelemetry, esses traces podendo ser visualizados com o Zipkin.

## Como utilizar
O aplicativo está totalmente "dockerizado", por isso vamos precisar apenas ter o Dokcer instalado. 

### Variáveis de ambiente

1. Input-Service:

- ORCHESTRATION_SERVICE_URL = http://orchestration-service:8081/
- SERVICE = input-service
- OTEL_URL = otel-collector:4317

2. Orchestration-service:

- WEATHER_KEY = {SUA_API_KEY} [WeatherAPI](https://www.weatherapi.com/)
- SERVICE = orchestration-service
- OTEL_URL = otel-collector:4317

### Executando aplicação

Para executar a aplicação basta rodar o seguinte comando:
```bash
docker-compose up
```
Podemos realizar uma request utilizando o `curl` ou qualquer cliente HTTP de sua preferência:
```bash
curl --request POST --url 'http://localhost:8080' -H "Content-Type: application/json" -d '{"cep" : "17560015"}'
```

Para acessar o Zipkin utilize o seguinte endereço:
```bash
http://127.0.0.1:9411
```
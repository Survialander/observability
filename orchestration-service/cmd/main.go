package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Survialander/orchestration-service/configs"
	"github.com/Survialander/orchestration-service/internal"
	"github.com/Survialander/orchestration-service/internal/handlers"
	"github.com/go-chi/chi"
)

func main() {
	err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}

	signChannel := make(chan os.Signal, 1)
	signal.Notify(signChannel, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := configs.InitOTel()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal("failed to shutdown tracer provider: ", err)
		}
	}()

	cepService := internal.NewCepService()
	weatherService := internal.NewWeatherService()
	weatherHandler := handlers.NewWeatherHandler(cepService, weatherService)

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		weatherHandler.GetWeather(w, r)
	})

	go func() {
		log.Println("Runing orchestration service on port 8081")
		if err := http.ListenAndServe(":8081", router); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-signChannel:
		log.Println("shutting down server gracefully...")
	case <-ctx.Done():
		cancel()
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}

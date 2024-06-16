package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Survialander/input-service/configs"
	"github.com/Survialander/input-service/internal/handler"
	"github.com/go-chi/chi"
)

func main() {
	configs.LoadConfig(".")

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

	r := chi.NewRouter()

	weatherInfoHandler := handler.NewWeatherHandler()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		weatherInfoHandler.GetWeatherInfo(w, r)
	})

	go func() {
		log.Println("Runing input service on port 8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
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

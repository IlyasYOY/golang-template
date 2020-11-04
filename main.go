package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	viper.AutomaticEnv()

	viper.SetDefault("request.timeout.duration", time.Second*2)
	viper.SetDefault("application.port", ":8080")

	requestTimeoutDuration := viper.GetDuration("request.timeout.duration")
	applicationPort := viper.GetString("application.port")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(requestTimeoutDuration))
	r.Use(middleware.ContentCharset("utf-8"))
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			log.Fatal("Error writing response.", err)
		}
	})

	log.Info("Starting application...")
	err := http.ListenAndServe(applicationPort, r)
	if err != nil {
		log.Fatal("Error starting application", err)
	}
}

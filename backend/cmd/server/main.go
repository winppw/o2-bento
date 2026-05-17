package main

import (
	"log"
	"net/http"
	"os"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"

	"github.com/o2ai/launch-assistant/backend/internal/handler"
	"github.com/o2ai/launch-assistant/backend/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS)
	// 60 requests per minute per IP — protects against polling abuse
	r.Use(httprate.LimitByIP(60, time.Minute))

	r.Get("/status", handler.GetStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

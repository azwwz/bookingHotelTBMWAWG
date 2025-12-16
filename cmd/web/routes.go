package main

import (
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	return r
}

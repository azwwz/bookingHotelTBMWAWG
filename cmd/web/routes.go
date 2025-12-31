package main

import (
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/pkg/config"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(NoSurf)
	r.Use(SessionLoad)
	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	r.Get("/generals", handlers.Repo.Generals)
	r.Get("/generals", handlers.Repo.Generals)
	r.Get("/majors", handlers.Repo.Majors)
	r.Get("/search-availability", handlers.Repo.Availability)
	r.Post("/search-availability", handlers.Repo.PostAvailability)
	r.Get("/make-reservation", handlers.Repo.Reservation)
	r.Get("/contact", handlers.Repo.Contact)

	//fileserver return handler get file
	FileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", FileServer))
	return r
}

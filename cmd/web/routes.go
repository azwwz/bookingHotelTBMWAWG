package main

import (
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/handlers"
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

	// get into search availability page
	r.Get("/search-availability", handlers.Repo.Availability)
	// post date and search availability rooms
	r.Post("/search-availability", handlers.Repo.PostAvailability)
	// in certain page get availability if true
	r.Post("/search-availability-json", handlers.Repo.AvailabilityJson)
	// after search-availability get a page and choose room
	r.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	// after search-availability-json we book now
	r.Get("/book-room", handlers.Repo.BookRoom)

	// get  make-reservation.page
	r.Get("/make-reservation", handlers.Repo.Reservation)

	// post reservation form and turn to reservation summery
	r.Post("/make-reservation", handlers.Repo.PostReservation)

	// after post /make-reservation get  reservation-summary page
	r.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	r.Get("/contact", handlers.Repo.Contact)

	// route to login page
	r.Get("/user/login", handlers.Repo.ShowLogin)
	//fileserver return handler get file
	FileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", FileServer))
	return r
}

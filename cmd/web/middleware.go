package main

import (
	"github.com/azwwz/bookingHotelTBMWAWG/internal/helpers"
	"net/http"

	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	noSurfHandler := nosurf.New(next)

	noSurfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   app.InProduction,
	})

	return noSurfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			sessionManager.Put(r.Context(), "error", "log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

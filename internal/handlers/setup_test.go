package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

var app *config.AppConfig
var sessionManager *scs.SessionManager
var pathToTemplate = "../.."
var functions template.FuncMap

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app = &config.AppConfig{}
	app.InProduction = false

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction

	app.SessionManager = sessionManager

	app.UseCache = true
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc

	Repo = NewRepo(app)

	render.NewTemplates(app)

	r := chi.NewRouter()
	// r.Use(NoSurf)
	r.Use(SessionLoad)
	r.Get("/", Repo.Home)
	r.Get("/about", Repo.About)
	r.Get("/generals", Repo.Generals)
	r.Get("/generals", Repo.Generals)
	r.Get("/majors", Repo.Majors)

	r.Get("/search-availability", Repo.Availability)
	r.Post("/search-availability", Repo.PostAvailability)
	r.Post("/search-availability-json", Repo.AvailabilityJson)

	r.Get("/make-reservation", Repo.Reservation)
	r.Post("/make-reservation", Repo.PostReservation)
	r.Get("/reservation-summary", Repo.ReservationSummary)

	r.Get("/contact", Repo.Contact)

	//fileserver return handler get file
	FileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", FileServer))
	return r
}
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

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/templates/*.page.html", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	_, err = filepath.Glob(fmt.Sprintf("%s/templates/*.layout.html", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/templates/*.layout.html", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}

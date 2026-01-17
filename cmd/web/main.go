package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/handlers"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
)

var sessionManager *scs.SessionManager
var app *config.AppConfig

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	// start the server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func run() error {
	gob.Register(models.Reservation{})

	// create golbal app config
	app = &config.AppConfig{}

	app.InProduction = false

	// session manager create bind to the app config
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction

	app.SessionManager = sessionManager

	// let handler use the app config
	repo := handlers.NewRepo(app)
	handlers.SetRepo(repo)

	// let render package  use the app config
	render.NewTemplates(app)

	// create template cache bind to the app config
	tc, err := render.CreateTemplateCache()
	if err != nil {
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	return nil
}

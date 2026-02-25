package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/driver"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/handlers"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository/dbrepo"
)

var sessionManager *scs.SessionManager
var app *config.AppConfig

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
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

func run() (*driver.DB, error) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestrictions{})

	// create golbal app config
	app = &config.AppConfig{}

	app.InProduction = false

	// config appConfig log.logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session manager create bind to the app config
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction

	app.SessionManager = sessionManager

	// new database connection
	log.Println("conecting to database ... ")
	dbconn, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=123456")
	if err != nil {
		return nil, err
	}

	dbrepo := dbrepo.NewPostgresDBRepo(dbconn.SQL, app)
	// let handler use the app config
	repo := handlers.NewRepo(app, dbrepo)
	handlers.SetRepo(repo)

	// let render package  use the app config
	render.NewRender(app)

	// create template cache bind to the app config
	tc, err := render.CreateTemplateCache()
	if err != nil {
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	return dbconn, nil
}

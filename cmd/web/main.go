package main

import (
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/pkg/config"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/handlers"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/render"
)

func main() {
	app := &config.AppConfig{}

	repo := handlers.NewRepo(app)
	handlers.SetRepo(repo)
	render.NewTemplates(app)
	tc, err := render.CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

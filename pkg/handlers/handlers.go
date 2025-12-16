package handlers

import (
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/pkg/config"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/models"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func SetRepo(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{})
}

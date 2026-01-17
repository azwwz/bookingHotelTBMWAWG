package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/forms"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
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
	remoteIP := r.RemoteAddr
	repo.App.SessionManager.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again. there is about page."

	remoteIP := repo.App.SessionManager.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Fprintf(w, "yout start time is %s, endtime is %s", start, end)
}

type jsonAvailabilityJson struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonAvailabilityJson{
		Ok:      true,
		Message: "2025年12月31日",
	}
	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(out))
	w.Header().Set("content-type", "application/json")
	w.Write(out)
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: &forms.Form{},
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	form := forms.NewForm(r.PostForm)
	form.Require("first_name", "last_name", "email", "phone")
	form.Minimum("first_name", 3)
	form.IsEmail("email")
	reversation := &models.Reservation{
		First_name: r.Form.Get("first_name"),
		Last_name:  r.Form.Get("last_name"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
	}
	if !form.Valid() {
		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}
	repo.App.SessionManager.Put(r.Context(), "reservation", reversation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get item form session")
		repo.App.SessionManager.Put(r.Context(), "error", "can not get item from session")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	repo.App.SessionManager.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}

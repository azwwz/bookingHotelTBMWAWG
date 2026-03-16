package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/forms"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository"
	"github.com/go-chi/chi/v5"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, d repository.DatabaseRepo) *Repository {
	return &Repository{
		App: a,
		DB:  d,
	}
}

func TestNewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.TestNewPostgresDBRepo(a),
	}
}

func SetRepo(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	// time : 01/02 03:04:05PM '06 -0700 Mon Jan
	layout := "2006-1-2"
	startData, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot  time.Parse(layout, sd)")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	endData, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot time.Parse(layout, ed)")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startData, endData)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot SearchAvailabilityForAllRooms")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if len(rooms) == 0 {
		// no availability
		repo.App.SessionManager.Put(r.Context(), "warning", "no availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startData,
		EndDate:   endData,
	}

	repo.App.SessionManager.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

type jsonAvailabilityJson struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    int    `json:"room_id"`
}

func (repo *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonAvailabilityJson{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	startDate, _ := time.Parse("2006-01-02", sd)
	endDate, _ := time.Parse("2006-01-02", ed)

	available, err := repo.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonAvailabilityJson{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return

	}
	resp := jsonAvailabilityJson{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    roomID,
	}

	// I removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "  ")

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(out)
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// now res has 	ID int；  StartDate time.Time; EndDate   time.Time
	res, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.SessionManager.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["StartData"] = sd
	stringMap["EndData"] = ed

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can't get Room By ID from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	// reservation has the all reservation and room name and id
	repo.App.SessionManager.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      &forms.Form{},
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can not get PostForm from request")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.NewForm(r.PostForm)
	form.Require("first_name", "last_name", "email", "phone")
	form.Minimum("first_name", 3)
	form.IsEmail("email")

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	// after check the valid, stroe into database

	// Test  repo.DB.InsertReservation(reservation) failed
	reservation.ID, err = repo.DB.InsertReservation(reservation)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can not InsertReservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservation.ID,
		RestrictionID: 1,
	}
	// test repo.DB.InsertRoomRestriction(restriction) failed
	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can not InsertRoomRestriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.App.SessionManager.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.SessionManager.Put(r.Context(), "error", "can not get item from session")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	repo.App.SessionManager.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom receive a room_id and return to /make-reservation page with reservation and room info
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "room_id"))
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "can not get URLParam")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot get a reservation form session")
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return

	}

	res.RoomID = roomID

	repo.App.SessionManager.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot get r.URL.Query")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	startDate := r.URL.Query().Get("s")

	endDate := r.URL.Query().Get("e")

	var res models.Reservation
	room, err := repo.DB.GetRoomByID(roomID)

	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot  GetRoomByID")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	res.StartDate, _ = time.Parse("2006-01-02", startDate)
	res.EndDate, _ = time.Parse("2006-01-02", endDate)
	res.Room = room
	res.RoomID = roomID

	repo.App.SessionManager.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

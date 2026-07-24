package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/helpers"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/forms"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/render"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository"
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

// SetRepo 设置仓库的全局变量
// 参数:
//
//	r *Repository: 要设置的仓库指针
func SetRepo(r *Repository) {
	Repo = r // 将传入的仓库指针赋值给全局变量Repo
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	// time : 01/02 03:04:05PM '06 -0700 Mon Jan
	layout := "2006-1-2"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot  time.Parse(layout, sd)")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot time.Parse(layout, ed)")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
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
		StartDate: startDate,
		EndDate:   endDate,
	}

	repo.App.SessionManager.Put(r.Context(), "reservation", res)

	err = render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot render")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
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
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
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
	stringMap["StartDate"] = sd
	stringMap["EndDate"] = ed

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

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      &forms.Form{},
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation post reservation form and turn to reservation summery
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

	// get roomID
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// from roomID get room
	room, err := repo.DB.GetRoomByID(roomID)
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "cannot get Room by ID")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	// add all to reservation for reservation summery
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		err := render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			repo.App.SessionManager.Put(r.Context(), "error", "error rendering template")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		return
	}
	// after check the valid, store into database

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

	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s :,<br>
		This is confirm your reservation from %s to %s .`,
		reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       "yue@yue.com",
		From:     "19704769@qq.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	repo.App.MailChan <- msg

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

	err := render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		repo.App.SessionManager.Put(r.Context(), "error", "error rendering template")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

// ChooseRoom receive a room_id and return to /make-reservation page with reservation and room info
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
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

func (repo *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.NewForm(nil),
	})
}

func (repo *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.SessionManager.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.NewForm(r.PostForm)
	form.Require("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := repo.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		repo.App.SessionManager.Put(r.Context(), "error", "authentication failed")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	repo.App.SessionManager.Put(r.Context(), "user_id", id)
	repo.App.SessionManager.Put(r.Context(), "flash", "logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.SessionManager.Destroy(r.Context())
	_ = repo.App.SessionManager.RenewToken(r.Context())
	repo.App.SessionManager.Put(r.Context(), "flash", "logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (repo *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows all new reservations in admin tool
func (repo *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllNewReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations shows all reservations in admin tool
func (repo *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllReservations()
	if err != nil {
		log.Println(err)
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) AdminNewReservation(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllReservations()
	if err != nil {
		log.Println(err)
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

// AdminShowReservation shows a reservation in admin tool
func (repo *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	// get id from url use splite
	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(err)
		return
	}

	src := parts[3]
	stringmap := make(map[string]string)
	stringmap["src"] = src

	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		log.Println(err)
	}

	// create data
	data := make(map[string]interface{})

	// add reservation to data
	data["reservation"] = reservation
	render.Template(w, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringmap,
		Form:      forms.NewForm(nil),
	})
}

// AdminPostShowReservation updates a reservation in admin tool
func (repo *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	// get id from url use split
	urlParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(urlParts[len(urlParts)-1])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form
	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get reservation by id
	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create map save src
	src := urlParts[len(urlParts)-2]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	// check form
	form := forms.NewForm(r.PostForm)
	form.Require("first_name", "last_name", "email", "phone")
	form.Minimum("first_name", 3)
	form.IsEmail("email")
	if !form.Valid() {
		reservation.FirstName = form.Get("first_name")
		reservation.LastName = form.Get("last_name")
		reservation.Email = form.Get("email")
		reservation.Phone = form.Get("phone")
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	reservation.FirstName = form.Get("first_name")
	reservation.LastName = form.Get("last_name")
	reservation.Email = form.Get("email")
	reservation.Phone = form.Get("phone")

	// update reservation
	err = repo.DB.UpdateReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// redirect to reservation show page
	repo.App.SessionManager.Put(r.Context(), "flash", "reservation updated successfully")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (repo *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	// get id from url use split
	urlParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(urlParts[len(urlParts)-1])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	repo.DB.UpdateProcessedForReservation(id, 1)
	src := urlParts[len(urlParts)-2]
	repo.App.SessionManager.Put(r.Context(), "flash", "Reservation processed successfully")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminDeleteReservation deletes a reservation in admin tool
func (repo *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	err = repo.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")
	repo.App.SessionManager.Put(r.Context(), "flash", "Reservation deleted successfully")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminReservationsCalendar shows all reservations in admin tool calendar view
func (repo *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	lastMonth := last.Format("01")

	nextMonthYear := next.Format("2006")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get first day and last day of the month ----------
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	// set day = 0 means the last day of the previous month
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 0, 0, 0, 0, currentLocation)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	// get all rooms ----------
	rooms, err := repo.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

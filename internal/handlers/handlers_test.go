package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "get", http.StatusOK},
	{"about", "/about", "get", http.StatusOK},
	{"generals", "/generals", "get", http.StatusOK},
	{"search-availability", "/search-availability", "get", http.StatusOK},
	{"majors", "/majors", "get", http.StatusOK},
	{"contact", "/contact", "post", http.StatusOK},
	//{"search-availability", "/search-availability", "post", []postData{
	//	{key: "start", value: "2012-02-10"},
	//	{key: "end", value: "2012-02-10"},
	//}, http.StatusOK},
	//{"search-availability-json", "/search-availability-json", "post", []postData{
	//	{key: "start", value: "2012-02-10"},
	//	{key: "end", value: "2012-02-10"},
	//}, http.StatusOK},
	//{"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},
	//{"make-reservation", "/make-reservation", "post", []postData{
	//	{key: "first_name", value: "555"},
	//	{key: "last_name", value: "5555"},
	//	{key: "email", value: "123@123.com"},
	//	{key: "phone", value: "123123"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	server := httptest.NewTLSServer(routes)
	defer server.Close()

	for _, e := range theTests {
		res, err := server.Client().Get(server.URL + e.url)
		if err != nil {
			t.Log(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
		}
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// Test Success
	values := url.Values{}
	values.Add("start", "2050-01-01")
	values.Add("end", "2050-01-02")

	request, _ := http.NewRequest("POST", "/availability", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := getCtx(request)
	request = request.WithContext(ctx)

	handlerFunc := http.HandlerFunc(Repo.PostAvailability)
	recorder := httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}

	// 2. startDate, err := time.Parse(layout, sd)
	values = url.Values{}
	values.Add("start", "invalid")
	values.Add("end", "2050-01-02")

	request, _ = http.NewRequest("POST", "/availability", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	handlerFunc = http.HandlerFunc(Repo.PostAvailability)
	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusMovedPermanently)
	}

	// 2. endDate, err := time.Parse(layout, sd)
	values = url.Values{}
	values.Add("start", "2050-01-01")
	values.Add("end", "invalid")

	request, _ = http.NewRequest("POST", "/availability", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusMovedPermanently)
	}

	// 3 . rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	values = url.Values{}
	values.Add("start", "2050-01-01")
	values.Add("end", "2050-01-03")

	request, _ = http.NewRequest("POST", "/availability", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusMovedPermanently)
	}

	// 4. if len(rooms) == 0
	values = url.Values{}
	values.Add("start", "2050-01-01")
	values.Add("end", "2050-01-04")

	request, _ = http.NewRequest("POST", "/availability", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_Reservation(t *testing.T) {

	// test through
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	request, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(request)
	request = request.WithContext(ctx)
	sessionManager.Put(request.Context(), "reservation", reservation)
	recorder := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(Repo.Reservation)
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusOK)
	}

	// test No Reservation in session
	requestNOreservation, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx2 := getCtx(requestNOreservation)
	requestNOreservation = request.WithContext(ctx2)

	rr2 := httptest.NewRecorder()
	handlerFunc2 := http.HandlerFunc(Repo.Reservation)
	handlerFunc2.ServeHTTP(rr2, requestNOreservation)
	if rr2.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusTemporaryRedirect)
	}

	// Test get Reservation RoomID failed
	reservation.RoomID = 100
	sessionManager.Put(request.Context(), "reservation", reservation)
	reservationM, ok := sessionManager.Get(requestNOreservation.Context(), "reservation").(models.Reservation)
	if ok {
		fmt.Println(reservationM.RoomID)
	}
	rr3 := httptest.NewRecorder()

	handlerFunc.ServeHTTP(rr3, request)
	if rr3.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// test pass
	reqBody := url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "1")

	r1, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 := getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 := httptest.NewRecorder()
	PostReservation := http.HandlerFunc(Repo.PostReservation)
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusSeeOther)
	}

	// test parseForm failed
	r2, _ := http.NewRequest("POST", "/make-reservation", nil)
	// request header
	ctx2 := getCtx(r1)
	r2 = r2.WithContext(ctx2)
	r2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr2 := httptest.NewRecorder()
	PostReservation.ServeHTTP(rr2, r2)
	if rr2.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}
	// Test start date failed

	reqBody = url.Values{}
	reqBody.Add("start_date", "invalid")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "1")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}
	// Test end date failed
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "invalid")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "1")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}
	// test room_id failed
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "invalid")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}

	// test invalid form
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "1")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusOK)
	}

	// test invalid InsertReservation
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "100")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}
	// test invalid InsertRoomRestriction
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "John")
	reqBody.Add("last_name", "Smith")
	reqBody.Add("email", "mail@mail.com")
	reqBody.Add("phone", "123123123")
	reqBody.Add("room_id", "2")

	r1, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	// request header
	ctx1 = getCtx(r1)
	r1 = r1.WithContext(ctx1)
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr1 = httptest.NewRecorder()
	PostReservation.ServeHTTP(rr1, r1)
	if rr1.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got -> %d, wanted-> %d",
			rr1.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJson(t *testing.T) {
	// first case - rooms are not available
	respBody := "start=2050-01-01"
	respBody = fmt.Sprintf("%s&%s", respBody, "end=2050-01-02")
	respBody = fmt.Sprintf("%s&%s", respBody, "room_id=1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(respBody))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// make handler handlerFunc
	handlerFunc := http.HandlerFunc(Repo.AvailabilityJson)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to out handler
	handlerFunc.ServeHTTP(rr, req)

	// var j jsonResponse
	var j jsonAvailabilityJson
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Errorf("show false but true")
	}

	// 	err := r.ParseForm()
	uv := url.Values{}
	uv.Add("start", "2050-01-01")
	uv.Add("end", "2050-01-02")
	uv.Add("room_id", "1")
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}
	if j.OK {
		t.Errorf("show true but false")
	}

	// test available, err := repo.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	uv = url.Values{}
	uv.Add("start", "2049-12-31")
	uv.Add("end", "2050-01-02")
	uv.Add("room_id", "1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(uv.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}
	if j.OK || j.Message != "Error querying database" {
		t.Errorf("show true but false")
	}

}

func TestRepository_ReservationSummary(t *testing.T) {
	request := httptest.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(request)
	request = request.WithContext(ctx)
	reservation := models.Reservation{
		StartDate: time.Time{},
		EndDate:   time.Time{},
	}
	app.SessionManager.Put(request.Context(), "reservation", reservation)
	handlerFunc := http.HandlerFunc(Repo.ReservationSummary)
	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, request)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", rr.Code, http.StatusOK)
	}
	// 1. reservation, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	request = httptest.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	rr = httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, request)
	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", rr.Code, http.StatusMovedPermanently)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	// test success
	request := httptest.NewRequest("GET", "/choose-room/1", nil)
	request = request.WithContext(getCtx(request))

	// add chi URLParam
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")
	chiCtx := context.WithValue(request.Context(), chi.RouteCtxKey, routeContext)
	request = request.WithContext(chiCtx)

	app.SessionManager.Put(request.Context(), "reservation", models.Reservation{})
	handlerFunc := http.HandlerFunc(Repo.ChooseRoom)
	rr := httptest.NewRecorder()
	request.RequestURI = "/choose-room/1"
	handlerFunc.ServeHTTP(rr, request)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", rr.Code, http.StatusSeeOther)
	}

	// test  roomID, err := LParam(r, "id"))
	request = httptest.NewRequest("GET", "/choose-room/1", nil)
	request = request.WithContext(getCtx(request))

	// add chi URLParam
	routeContext = chi.NewRouteContext()
	routeContext.URLParams.Add("id", "")
	chiCtx = context.WithValue(request.Context(), chi.RouteCtxKey, routeContext)
	request = request.WithContext(chiCtx)

	app.SessionManager.Put(request.Context(), "reservation", models.Reservation{})
	rr = httptest.NewRecorder()
	request.RequestURI = "/choose-room/1"
	handlerFunc.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test res, ok := repo.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	request = httptest.NewRequest("GET", "/choose-room/1", nil)
	request = request.WithContext(getCtx(request))

	// add chi URLParam
	routeContext = chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")
	chiCtx = context.WithValue(request.Context(), chi.RouteCtxKey, routeContext)
	request = request.WithContext(chiCtx)

	app.SessionManager.Put(request.Context(), "123", models.Reservation{})
	rr = httptest.NewRecorder()
	request.RequestURI = "/choose-room/1"
	handlerFunc.ServeHTTP(rr, request)
	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", rr.Code, http.StatusMovedPermanently)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	// test success
	request := httptest.NewRequest("GET", "/book-room?id=1&s=2050-01-01&e=2050-01-02", nil)
	ctx := getCtx(request)
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(Repo.BookRoom)
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusSeeOther {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusSeeOther)
	}

	// roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	request = httptest.NewRequest("GET", "/book-room?id=&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusMovedPermanently)
	}

	// room, err := repo.DB.GetRoomByID(roomID)
	request = httptest.NewRequest("GET", "/book-room?id=3&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	recorder = httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong response code: got -> %d, wanted-> %d", recorder.Code, http.StatusMovedPermanently)
	}

}
func getCtx(r *http.Request) context.Context {
	ctx, _ := sessionManager.Load(r.Context(), r.Header.Get("X-Session any word"))
	return ctx
}

package handlers

import (
	"context"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "get", []postData{}, http.StatusOK},
	{"about", "/about", "get", []postData{}, http.StatusOK},
	{"generals", "/generals", "get", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "get", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "post", []postData{
		{key: "start", value: "2012-02-10"},
		{key: "end", value: "2012-02-10"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "post", []postData{
		{key: "start", value: "2012-02-10"},
		{key: "end", value: "2012-02-10"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "post", []postData{
		{key: "first_name", value: "555"},
		{key: "last_name", value: "5555"},
		{key: "email", value: "123@123.com"},
		{key: "phone", value: "123123"},
	}, http.StatusOK},
}

func TestHanlders(t *testing.T) {
	routes := getRoutes()
	server := httptest.NewTLSServer(routes)
	defer server.Close()

	for _, e := range theTests {
		if e.method == "get" {
			res, err := server.Client().Get(server.URL + e.url)
			if err != nil {
				t.Log(err)
			}
			if res.StatusCode != http.StatusOK {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, p := range e.params {
				values.Add(p.key, p.value)
				log.Printf(" values.Add(%v,%v)\n ", p.key, p.value)
			}
			resp, err := server.Client().PostForm(server.URL+e.url, values)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
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
}

func getCtx(r *http.Request) context.Context {
	ctx, _ := sessionManager.Load(r.Context(), r.Header.Get("X-Session any word"))

	return ctx
}

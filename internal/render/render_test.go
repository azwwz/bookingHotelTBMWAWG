package render

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
)

func TestAddDefaultTemplate(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	sessionManager.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "../../templates/"

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	mywriter := myWriter{}
	err = RenderTemplate(&mywriter, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	err = RenderTemplate(&mywriter, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r := httptest.NewRequest("get", "/some", nil)

	ctx := r.Context()
	log.Println(r.Header.Get("X-Session"))
	ctx, err := sessionManager.Load(ctx, r.Header.Get("X-Session"))
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)
	// log.Println("r.Header.Get(X-Session)--" + r.Header.Get("X-Session"))
	// log.Println("r.contest ", r.Context())
	return r, nil
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplate = "../../templates/"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"formatDate": FormatDate,
}

var app *config.AppConfig
var tc map[string]*template.Template
var pathToTemplate = "./templates"

func NewRender(a *config.AppConfig) {
	app = a
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatData
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.SessionManager.PopString(r.Context(), "flash")
	td.Warning = app.SessionManager.PopString(r.Context(), "warning")
	td.Error = app.SessionManager.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.SessionManager.Exists(r.Context(), "user_id") {
		td.ISAuthenticated = 1
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var err error
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			return err
		}
	}
	t, ok := tc[tmpl]
	if !ok {
		err = fmt.Errorf("can not get template -- %s", tmpl)
		return err
	}
	buffer := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err = t.Execute(buffer, td)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// CreateTemplateCache parses the templates once
// store them in cache map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	_, err = filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}

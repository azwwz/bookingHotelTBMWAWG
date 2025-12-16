package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/azwwz/bookingHotelTBMWAWG/pkg/config"
	"github.com/azwwz/bookingHotelTBMWAWG/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a

}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("can not get template", tmpl)
	}
	buffer := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buffer, td)
	if err != nil {
		log.Println(err)
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

// CreateTemplateCache parses the templates once
// store them in cache map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}
	_, err = filepath.Glob("./templates/*.layout.html")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		log.Println("name is currently ", name)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		ts, err = ts.ParseGlob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}

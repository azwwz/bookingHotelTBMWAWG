package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration settings.
type AppConfig struct {
	TemplateCache  map[string]*template.Template
	UseCache       bool
	SessionManager *scs.SessionManager
	InProduction   bool
	CSRFToken      string
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
}

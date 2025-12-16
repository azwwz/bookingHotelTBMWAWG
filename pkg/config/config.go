package config

import "html/template"

// AppConfig holds the application configuration settings.
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
}

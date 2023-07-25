package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/brattonross/website/internal/theme"
)

//go:embed html/*
var templates embed.FS
var layoutPath = "html/layout.html"

func ParseTemplates(files ...string) (*template.Template, error) {
	if os.Getenv("DEV") == "true" {
		return template.ParseFiles(files...)
	}
	return template.ParseFS(templates, files...)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	theme := theme.GetTheme(r)
	if theme != "" {
		data["Theme"] = theme
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

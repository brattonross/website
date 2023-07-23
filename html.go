package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed html/*
var templates embed.FS
var layoutPath = "html/layout.html"

func ParseTemplates(files ...string) (*template.Template, error) {
	// TODO: Load from real file system in dev mode
	return template.ParseFS(templates, files...)
}

var themeCookieName = "brattonross_theme"

func SetTheme(w http.ResponseWriter, theme string) {
	cookie := http.Cookie{
		Name:     themeCookieName,
		Value:    theme,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   60 * 60 * 24 * 365 * 100, // 100 years
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	cookie, _ := r.Cookie(themeCookieName)
	if cookie != nil {
		data["Theme"] = cookie.Value
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

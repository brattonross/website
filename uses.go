package main

import (
	"net/http"
)

func UsesPage(w http.ResponseWriter) {
	tmpl, err := ParseTemplates(layoutPath, "html/uses.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", struct {
		Description string
		Title       string
	}{
		Description: "Stuff that I use on a daily basis.",
		Title:       "Uses",
	})
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

package main

import (
	"net/http"
)

func HomePage(w http.ResponseWriter) {
	tmpl, err := ParseTemplates(layoutPath, "html/index.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

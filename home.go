package main

import (
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := ParseTemplates(layoutPath, "html/index.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = RenderTemplate(w, r, tmpl, nil)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

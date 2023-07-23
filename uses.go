package main

import (
	"log"
	"net/http"
)

func UsesPage(w http.ResponseWriter) {
	tmpl, err := ParseTemplates(layoutPath, "html/uses.html")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

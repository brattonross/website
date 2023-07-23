package main

import (
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter) {
	tmpl, err := ParseTemplates(layoutPath, "html/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

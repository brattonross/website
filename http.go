package main

import (
	"log"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := ParseTemplates(layoutPath, "html/404.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = RenderTemplate(w, r, tmpl, map[string]interface{}{
		"Description": "Page not found",
		"Title":       "Not Found",
	})
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

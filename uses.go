package main

import (
	"net/http"
)

func UsesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := ParseTemplates(layoutPath, "html/uses.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = RenderTemplate(w, r, tmpl, map[string]interface{}{
		"Description": "Stuff that I use on a daily basis.",
		"Title":       "Uses",
	})
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

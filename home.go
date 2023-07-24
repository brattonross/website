package main

import (
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, posts []Post) {
	tmpl, err := ParseTemplates(layoutPath, "html/index.html")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = RenderTemplate(w, r, tmpl, map[string]interface{}{
		"Posts": posts,
	})
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

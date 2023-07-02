package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed html/*
var templates embed.FS
var layoutPath = "html/layout.html"

func parseTemplates(files ...string) (*template.Template, error) {
	// TODO: Load from real file system in dev mode
	return template.ParseFS(templates, files...)
}

func embeddedFileExists(path string) bool {
	_, err := templates.Open(path)
	return err == nil
}

func PageFromPath(w http.ResponseWriter, path string) {
	if strings.HasSuffix(path, "/") {
		path = path + "index.html"
	} else if !strings.HasSuffix(path, ".html") {
		path = path + ".html"
	}

	path = "html" + path
	log.Printf("Serving %s\n", path)

	if !embeddedFileExists(path) {
		w.WriteHeader(http.StatusNotFound)
		tmpl, err := parseTemplates(layoutPath, "html/404.html")
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
		return
	}

	tmpl, err := parseTemplates(layoutPath, path)
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

package main

import (
	"embed"
	"html/template"
)

//go:embed html/*
var templates embed.FS
var layoutPath = "html/layout.html"

func ParseTemplates(files ...string) (*template.Template, error) {
	// TODO: Load from real file system in dev mode
	return template.ParseFS(templates, files...)
}

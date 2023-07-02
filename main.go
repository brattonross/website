package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed all:public
var public embed.FS

func withCaching(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}

func main() {
	public, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/public/", withCaching(http.StripPrefix("/public/", http.FileServer(http.FS(public)))))

	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		RootPage(w)
	})

	http.HandleFunc("/blog/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/blog/")
		PostPage(w, slug)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		PageFromPath(w, r.URL.Path)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

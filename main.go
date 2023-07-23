package main

import (
	"embed"
	"encoding/xml"
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

	http.HandleFunc("/blog.json", func(w http.ResponseWriter, r *http.Request) {
		feed, err := GenerateBlogFeed()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		json, err := feed.ToJSON()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(json))
		if err != nil {
			InternalServerError(w, err)
			return
		}
	})

	http.HandleFunc("/blog/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		feed, err := GenerateBlogFeed()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		rss, err := feed.ToRss()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		_, err = w.Write([]byte(rss))
		if err != nil {
			InternalServerError(w, err)
			return
		}
	})

	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		type Sitemap struct {
			XMLName xml.Name `xml:"urlset"`
			Xmlns   string   `xml:"xmlns,attr"`
			Xsi     string   `xml:"xmlns:xsi,attr"`
			Schema  string   `xml:"xsi:schemaLocation,attr"`
			Urls    []struct {
				Loc string `xml:"loc"`
			} `xml:"url"`
		}

		sitemap := Sitemap{
			Xmlns:  "http://www.sitemaps.org/schemas/sitemap/0.9",
			Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
			Schema: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
		}

		urls := []string{
			"https://brattonross.xyz",
			"https://brattonross.xyz/uses",
			"https://brattonross.xyz/blog",
		}

		posts, err := listPosts()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		for _, post := range posts {
			urls = append(urls, "https://brattonross.xyz/blog/"+post.Slug)
		}

		for _, url := range urls {
			sitemap.Urls = append(sitemap.Urls, struct {
				Loc string `xml:"loc"`
			}{url})
		}

		output, err := xml.MarshalIndent(sitemap, "  ", "    ")
		if err != nil {
			InternalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		_, err = w.Write([]byte(xml.Header + string(output)))
		if err != nil {
			InternalServerError(w, err)
			return
		}
	})

	http.HandleFunc("/uses", func(w http.ResponseWriter, r *http.Request) {
		UsesPage(w)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			NotFound(w)
			return
		}

		HomePage(w)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

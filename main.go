package main

import (
	"embed"
	"encoding/xml"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/brattonross/website/internal/blog"
	"github.com/brattonross/website/internal/theme"
	gomd "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/feeds"
)

var isDev = os.Getenv("DEV") == "true"

//go:embed all:public
var public embed.FS

//go:embed data/blog/*.md
var posts embed.FS

//go:embed html/*
var templates embed.FS

type devBlogFS struct{}

func (f *devBlogFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join("data", "blog", name))
}

func (f *devBlogFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join("data", "blog", name))
}

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

	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		if isDev {
			http.StripPrefix("/public/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
			return
		}

		withCaching(http.StripPrefix("/public/", http.FileServer(http.FS(public)))).ServeHTTP(w, r)
	})

	var blogFS *blog.FS
	if isDev {
		blogFS = blog.NewFS(&devBlogFS{})
	} else {
		subs, err := fs.Sub(posts, "data/blog")
		if err != nil {
			log.Fatal(err)
		}
		blogFS = blog.NewFS(subs.(fs.ReadDirFS))
	}

	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		filePath := "html/blog.html"
		posts, err := blogFS.ReadDir()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		frontmatters := []blog.PostFrontmatter{}
		for _, post := range posts {
			frontmatters = append(frontmatters, post.Frontmatter)
		}

		tmpl, err := ParseTemplates(layoutPath, filePath)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		err = RenderTemplate(w, r, tmpl, map[string]interface{}{
			"Description": "Ross Bratton's blog",
			"Posts":       frontmatters,
			"Title":       "Blog",
		})
		if err != nil {
			InternalServerError(w, err)
			return
		}
	})

	http.HandleFunc("/blog/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/blog/")
		post, err := blogFS.Open(slug + ".md")
		if err != nil {
			if os.IsNotExist(err) {
				NotFound(w, r)
				return
			}

			InternalServerError(w, err)
			return
		}

		parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
		doc := parser.Parse(post.Content)

		renderer := html.NewRenderer(html.RendererOptions{
			Flags: html.CommonFlags | html.HrefTargetBlank,
		})

		bs := gomd.Render(doc, renderer)

		tmpl, err := ParseTemplates(
			layoutPath,
			"html/blogpost.html",
		)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		tmpl, err = tmpl.Parse("{{define \"content\"}}" + string(bs) + "{{end}}")
		if err != nil {
			InternalServerError(w, err)
			return
		}

		err = RenderTemplate(w, r, tmpl, map[string]interface{}{
			"Description": post.Frontmatter.Description,
			"Post":        post.Frontmatter,
			"Title":       post.Frontmatter.Title,
		})
		if err != nil {
			InternalServerError(w, err)
			return
		}
	})

	generateBlogFeed := func() (*feeds.Feed, error) {
		posts, err := blogFS.ReadDir()
		if err != nil {
			return nil, err
		}

		feed := &feeds.Feed{
			Title:       "Ross Bratton - Blog",
			Link:        &feeds.Link{Href: "https://brattonross.xyz/blog"},
			Description: "Ross Bratton's blog",
			Author:      &feeds.Author{Name: "Ross Bratton", Email: "bratton.ross@gmail.com"},
			Created:     posts[0].Frontmatter.Date,
		}

		for _, post := range posts {
			feed.Items = append(feed.Items, &feeds.Item{
				Id:          post.Slug,
				Title:       post.Frontmatter.Title,
				Link:        &feeds.Link{Href: "https://brattonross.xyz" + post.Frontmatter.Href},
				Description: post.Frontmatter.Description,
				Created:     post.Frontmatter.Date,
			})
		}

		return feed, nil
	}

	http.HandleFunc("/blog.json", func(w http.ResponseWriter, r *http.Request) {
		feed, err := generateBlogFeed()
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
		feed, err := generateBlogFeed()
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

		posts, err := blogFS.ReadDir()
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

	http.HandleFunc("/theme", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		newTheme := r.Form.Get("theme")
		if newTheme == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theme.SetTheme(w, newTheme)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	})

	http.HandleFunc("/uses", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			NotFound(w, r)
			return
		}

		posts, err := blogFS.ReadDir()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		if len(posts) > 5 {
			posts = posts[:5]
		}

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
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

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

var layoutPath = "html/layout.html"

func ParseTemplates(files ...string) (*template.Template, error) {
	if isDev {
		return template.ParseFiles(files...)
	}
	return template.ParseFS(templates, files...)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	theme := theme.GetTheme(r)
	if theme != "" {
		data["Theme"] = theme
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

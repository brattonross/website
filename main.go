package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("GET /", http.FileServer(http.Dir("public")))
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/index.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("GET /blog", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir("content/blog")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type Article struct {
			Title       string
			Description string
			Slug        string
			Date        time.Time
		}

		articles := make([]Article, 0, len(files))
		for _, file := range files {
			article, err := os.ReadFile("content/blog/" + file.Name())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			frontmatter := make(map[string]string)
			for i, line := range strings.Split(string(article), "\n") {
				if i == 0 {
					continue
				}
				if line == "---" {
					break
				}
				parts := strings.SplitN(line, ":", 2)
				frontmatter[parts[0]] = strings.TrimSpace(parts[1])
			}

			date, err := time.Parse("2006-01-02", frontmatter["date"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			articles = append(articles, Article{
				Title:       frontmatter["title"],
				Description: frontmatter["description"],
				Date:        date,
				Slug:        strings.TrimSuffix(file.Name(), ".md"),
			})
		}

		data := struct {
			Articles []Article
		}{
			Articles: articles,
		}

		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/blog.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("GET /blog/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		type Article struct {
			Title   string
			Date    time.Time
			Content string
		}

		article, err := os.ReadFile("content/blog/" + slug + ".md")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		frontmatter := make(map[string]string)
		i := -1
		for _, line := range strings.Split(string(article), "\n") {
			i++
			if i == 0 {
				continue
			}

			if line == "---" {
				break
			}

			parts := strings.SplitN(line, ":", 2)
			frontmatter[parts[0]] = strings.TrimSpace(parts[1])
		}

		content := strings.Join(strings.Split(string(article), "\n")[i+1:], "\n")
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse([]byte(content))

		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		content = string(markdown.Render(doc, renderer))

		date, err := time.Parse("2006-01-02", frontmatter["date"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Article Article
		}{
			Article: Article{
				Title:   frontmatter["title"],
				Date:    date,
				Content: content,
			},
		}

		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/blog.$slug.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("GET /uses", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/uses.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server listening at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

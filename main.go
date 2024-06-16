package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:generate bun run build
func main() {
	if err := run(os.Getenv); err != nil {
		log.Fatal(err)
	}
}

func run(getenv func(string) string) error {
	host := getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := slog.Default()
	s := NewServer(logger)

	fmt.Printf("Server listening at http://%s:%s\n", host, port)
	return http.ListenAndServe(host+":"+port, s)
}

func NewServer(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("public")))
	mux.HandleFunc("GET /{$}", handleHomePage(logger))
	mux.HandleFunc("GET /blog", handleBlogPage(logger))
	mux.HandleFunc("GET /blog/{slug}", handleBlogPostPage(logger))
	mux.HandleFunc("GET /uses", handleUsesPage(logger))

	return mux
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func handle(fn handlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			logger.Error("http request error", "method", r.Method, "path", r.URL.Path, "error", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func handleHomePage(logger *slog.Logger) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/index.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, nil)
	}, logger)
}

func handleBlogPage(logger *slog.Logger) http.HandlerFunc {
	type article struct {
		Title       string
		Description string
		Slug        string
		Date        time.Time
	}
	type data struct {
		Articles []article
	}
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		files, err := os.ReadDir("content/blog")
		if err != nil {
			return err
		}

		as := make([]article, 0, len(files))
		for _, file := range files {
			a, err := os.ReadFile("content/blog/" + file.Name())
			if err != nil {
				return err
			}

			frontmatter := make(map[string]string)
			for i, line := range strings.Split(string(a), "\n") {
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
				return err
			}
			as = append(as, article{
				Title:       frontmatter["title"],
				Description: frontmatter["description"],
				Date:        date,
				Slug:        strings.TrimSuffix(file.Name(), ".md"),
			})
		}

		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/blog.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data{Articles: as})
	}, logger)
}

func handleBlogPostPage(logger *slog.Logger) http.HandlerFunc {
	type article struct {
		Title   string
		Date    time.Time
		Content string
	}
	type data struct {
		Article article
	}
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		slug := r.PathValue("slug")
		a, err := os.ReadFile("content/blog/" + slug + ".md")
		if err != nil {
			return err
		}

		frontmatter := make(map[string]string)
		i := -1
		for _, line := range strings.Split(string(a), "\n") {
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

		content := strings.Join(strings.Split(string(a), "\n")[i+1:], "\n")
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse([]byte(content))

		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		content = string(markdown.Render(doc, renderer))

		date, err := time.Parse("2006-01-02", frontmatter["date"])
		if err != nil {
			return err
		}

		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/blog.$slug.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data{
			Article: article{
				Title:   frontmatter["title"],
				Date:    date,
				Content: content,
			},
		})
	}, logger)
}

func handleUsesPage(logger *slog.Logger) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/uses.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, nil)
	}, logger)
}

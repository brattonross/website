package main

import (
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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	env := parseEnv()
	s := newServer()

	http.Handle("GET /", http.FileServer(http.Dir("public")))
	http.HandleFunc("GET /{$}", s.handleHomePage())
	http.HandleFunc("GET /blog", s.handleBlogPage())
	http.HandleFunc("GET /blog/{slug}", s.handleBlogPostPage())
	http.HandleFunc("GET /uses", s.handleUsesPage())

	s.logger.Info("Server listening at http://localhost:" + env.port)
	return http.ListenAndServe("0.0.0.0:"+env.port, nil)
}

type env struct {
	port string
}

func parseEnv() *env {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &env{port}
}

type server struct {
	logger *slog.Logger
}

func newServer() *server {
	return &server{
		logger: slog.Default(),
	}
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func (s *server) handle(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
	}
}

func (s *server) handleHomePage() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) error {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/index.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, nil)
	})
}

func (s *server) handleBlogPage() http.HandlerFunc {
	type article struct {
		Title       string
		Description string
		Slug        string
		Date        time.Time
	}
	type data struct {
		Articles []article
	}
	return s.handle(func(w http.ResponseWriter, r *http.Request) error {
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
	})
}

func (s *server) handleBlogPostPage() http.HandlerFunc {
	type article struct {
		Title   string
		Date    time.Time
		Content string
	}
	type data struct {
		Article article
	}
	return s.handle(func(w http.ResponseWriter, r *http.Request) error {
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
	})
}

func (s *server) handleUsesPage() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) error {
		tmpl, err := template.ParseFiles("templates/root.tmpl", "templates/uses.tmpl")
		if err != nil {
			return err
		}
		return tmpl.Execute(w, nil)
	})
}

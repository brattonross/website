package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/brattonross/website/markdown"
	gomd "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed data/blog/*.md
var posts embed.FS

// postFrontmatter represents the frontmatter of a blog post.
type postFrontmatter struct {
	Title       string
	Href        string
	Date        time.Time
	Description string
}

type post struct {
	Slug        string
	Frontmatter postFrontmatter
	Content     []byte
}

// parsePost parses a blog post from the given reader.
func parsePost(slug string, r io.Reader) (*post, error) {
	md, err := markdown.Parse(r)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", md.Frontmatter["date"])
	if err != nil {
		return nil, err
	}

	return &post{
		Slug: slug,
		Frontmatter: postFrontmatter{
			Title:       md.Frontmatter["title"],
			Href:        "/blog/" + slug,
			Date:        date,
			Description: md.Frontmatter["description"],
		},
		Content: md.Content,
	}, nil
}

func postByFileName(filename string) (*post, error) {
	file, err := posts.Open(filepath.Join("data", "blog", filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	slug := strings.TrimSuffix(filename, filepath.Ext(filename))
	post, err := parsePost(slug, file)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func listPosts() ([]post, error) {
	entries, err := posts.ReadDir(filepath.Join("data", "blog"))
	if err != nil {
		return nil, err
	}

	posts := []post{}
	for _, entry := range entries {
		post, err := postByFileName(entry.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}

func RootPage(w http.ResponseWriter) {
	filePath := "html/blog.html"
	posts, err := listPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	frontmatters := []postFrontmatter{}
	for _, post := range posts {
		frontmatters = append(frontmatters, post.Frontmatter)
	}

	tmpl, err := ParseTemplates(layoutPath, filePath)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", struct {
		Description string
		Posts       []postFrontmatter
		Title       string
	}{
		Description: "Ross Bratton's blog",
		Posts:       frontmatters,
		Title:       "Blog",
	})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func PostPage(w http.ResponseWriter, slug string) {
	post, err := postByFileName(slug + ".md")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.Parse("{{define \"content\"}}" + string(bs) + "{{end}}")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", struct {
		Description string
		Post        postFrontmatter
		Title       string
	}{
		Description: post.Frontmatter.Description,
		Post:        post.Frontmatter,
		Title:       post.Frontmatter.Title,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

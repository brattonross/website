package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/brattonross/website/internal/markdown"
)

//go:embed data/blog/*.md
var posts embed.FS

func openPost(filename string) (io.ReadCloser, error) {
	if os.Getenv("DEV") == "true" {
		return os.Open(filepath.Join("data", "blog", filename))
	}
	return posts.Open(filepath.Join("data", "blog", filename))
}

func readPostsDir() ([]fs.DirEntry, error) {
	if os.Getenv("DEV") == "true" {
		return os.ReadDir(filepath.Join("data", "blog"))
	}
	return posts.ReadDir(filepath.Join("data", "blog"))
}

// PostFrontmatter represents the frontmatter of a blog post.
type PostFrontmatter struct {
	Title       string
	Href        string
	Date        time.Time
	Description string
}

type Post struct {
	Slug        string
	Frontmatter PostFrontmatter
	Content     []byte
}

// parsePost parses a blog post from the given reader.
func parsePost(slug string, r io.Reader) (*Post, error) {
	md, err := markdown.Parse(r)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", md.Frontmatter["date"])
	if err != nil {
		return nil, err
	}

	return &Post{
		Slug: slug,
		Frontmatter: PostFrontmatter{
			Title:       md.Frontmatter["title"],
			Href:        "/blog/" + slug,
			Date:        date,
			Description: md.Frontmatter["description"],
		},
		Content: md.Content,
	}, nil
}

func PostByFilename(filename string) (*Post, error) {
	file, err := openPost(filename)
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

func ListPosts() ([]Post, error) {
	entries, err := readPostsDir()
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, entry := range entries {
		post, err := PostByFilename(entry.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Frontmatter.Date.After(posts[j].Frontmatter.Date)
	})

	return posts, nil
}

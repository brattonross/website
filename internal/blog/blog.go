package blog

import (
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/brattonross/website/internal/markdown"
)

// PostFrontmatter represents the frontmatter of a blog post.
type PostFrontmatter struct {
	Title       string
	Href        string
	Date        time.Time
	Description string
}

// Post represents a blog post.
type Post struct {
	Slug        string
	Frontmatter PostFrontmatter
	Content     []byte
}

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

// FS represents a blog using a filesystem.
type FS struct {
	fs fs.ReadDirFS
}

// NewFS creates a new blog using the given filesystem.
func NewFS(fs fs.ReadDirFS) *FS {
	return &FS{fs: fs}
}

// Open opens the blog post with the given filename.
func (fs *FS) Open(filename string) (*Post, error) {
	file, err := fs.fs.Open(filename)
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

// ReadDir reads the blog posts from the filesystem.
func (fs *FS) ReadDir() ([]Post, error) {
	entries, err := fs.fs.ReadDir(".")
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, entry := range entries {
		post, err := fs.Open(entry.Name())
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

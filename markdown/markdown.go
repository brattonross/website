package markdown

import (
	"bufio"
	"io"
	"strings"
)

// Frontmatter lives at the top of a markdown file, delimited by "---".
type Frontmatter map[string]string

// Markdown represents some markdown content.
type Markdown struct {
	Content     []byte
	Frontmatter Frontmatter
}

// Parse parses the markdown content from the given reader.
func Parse(r io.Reader) (*Markdown, error) {
	frontmatter := map[string]string{}
	content := []byte{}
	scanner := bufio.NewScanner(r)

	stage := "initial"

	for scanner.Scan() {
		if stage == "content" {
			content = append(content, scanner.Bytes()...)
			content = append(content, '\n')
			continue
		}

		line := scanner.Text()

		if line == "---" {
			if stage == "initial" {
				stage = "frontmatter"
				continue
			} else if stage == "frontmatter" {
				stage = "content"
				continue
			}
		}

		if line != "" && stage == "initial" {
			stage = "content"
		}

		if stage == "frontmatter" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			frontmatter[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			continue
		}
	}

	return &Markdown{
		Content:     content,
		Frontmatter: frontmatter,
	}, nil
}

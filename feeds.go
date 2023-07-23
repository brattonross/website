package main

import (
	"github.com/gorilla/feeds"
)

func GenerateBlogFeed() (*feeds.Feed, error) {
	posts, err := listPosts()
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

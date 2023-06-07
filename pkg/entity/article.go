// Package entity implements all state components
package entity

import "gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"

// Resource is a resource list for an endpoint.
type Resource struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous any       `json:"previous"`
	Results  []Article `json:"results"`
}

// Result is a resource list result.
type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	User  *User  `json:"user,omitempty"`
	Tags  []Tag  `json:"tags"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RequestGetArticles struct {
	Title  *string
	UserID *int
	Limit  int `validate:"gte=0,default=10"`
	Page   int `validate:"gte=0,default=1"`
}

type ResponseGetArticles struct {
	Items    []*Article
	Metadata rest.Pagination
}

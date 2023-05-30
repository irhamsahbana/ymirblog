// Package entity implements all state components
package entity

// Resource is a resource list for an endpoint.
type Resource struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous any       `json:"previous"`
	Results  []Article `json:"results"`
}

// Result is a resource list result.
type Article struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type RequestGetArticles struct {
	Name  *string
	Limit int `validate:"gte=0,default=10"`
	Page  int `validate:"gte=0,default=1"`
}

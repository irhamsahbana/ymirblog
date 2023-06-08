package rest

import "gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"

type ResponseArticle struct {
	Message  string           `json:"message"`
	Article  *entity.Article  `json:"article,omitempty"`
	Articles []entity.Article `json:"articles,omitempty"`
}

type RequestUpdateArticle struct {
	ID     string   `json:"id" validate:"required"`
	UserID int      `json:"user_id" validate:"required"`
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Tags   []string `json:"tags,omitempty" validate:"required"`
}

type RequestCreateArticle struct {
	UserID int      `json:"user_id" validate:"required"`
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Tags   []string `json:"tags,omitempty" validate:"required"`
}

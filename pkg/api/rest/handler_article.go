// Package rest is port handler.
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase/article"
)

// Article handler instance data.
type Article struct {
	UcArticle article.T
	DB        *ymirblog.Database
}

// Register is endpoint group for handler.
func (a *Article) Register(router chi.Router) {
	router.Get("/articles", rest.HandlerAdapter(a.GetAllArticle).JSON)
}

// GetAllArticle Handler
func (a *Article) GetAllArticle(w http.ResponseWriter, r *http.Request) ([]*entity.Article, error) {
	var (
		request entity.RequestGetArticles
	)
	b, err := rest.Bind(r, &request)
	if err != nil {
		return []*entity.Article{}, rest.ErrBadRequest(w, r, err)
	}
	if err := b.Validate(); err != nil {
		return []*entity.Article{}, rest.ErrBadRequest(w, r, err)
	}

	res, err := a.UcArticle.GetAll(r.Context(), r, request)
	if err != nil {
		return []*entity.Article{}, err
	}
	return res, nil
}

// Package rest is port handler.
package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase/article"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase/user"
)

// Article handler instance data.
type Article struct {
	UcArticle article.T
	UcUser    user.T
}

// Register is endpoint group for handler.
func (a *Article) Register(router chi.Router) {
	router.Get("/articles", rest.HandlerAdapter(a.GetAllArticle).JSON)
	router.Post("/articles", rest.HandlerAdapter(a.CreateArticle).JSON)
	router.Get("/articles/{id}", rest.HandlerAdapter(a.GetByIDArticle).JSON)
	router.Delete("/articles/{id}", rest.HandlerAdapter(a.DeleteArticle).JSON)
	router.Patch("/articles/{id}", rest.HandlerAdapter(a.UpdateArticle).JSON)
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

	res, err := a.UcArticle.GetAll(r.Context(), request)
	if err != nil {
		return []*entity.Article{}, err
	}

	rest.Paging(r, rest.Pagination{
		Page:  res.Metadata.Page,
		Limit: res.Metadata.Limit,
		Total: res.Metadata.Total,
	})

	return res.Items, nil
}

// CreateArticle Handler
func (a *Article) CreateArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	// binding and validate request body
	request := RequestCreateArticle{}
	b, err := rest.Bind(r, &request)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}
	if err = b.Validate(); err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping request to entity
	// tags
	tags := []entity.Tag{}
	for _, tag := range request.Tags {
		t := entity.Tag{
			Name: tag,
		}

		tags = append(tags, t)
	}
	// user
	u, err := a.UcUser.GetUserID(r.Context(), request.UserID)
	if err != nil {
		return res, rest.ErrNotFound(w, r, err)
	}

	entity := entity.Article{
		Title: request.Title,
		Body:  request.Body,
		Tags:  tags,
		User:  &u,
	}

	//create entity with usecase
	ent, err := a.UcArticle.Create(r.Context(), entity)
	if err != nil {
		log.Error().Err(err).Msg("CreateArticle4")
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping entity to response
	res.Message = "success create article"
	res.Article = &ent

	return res, nil
}

// Delete Article handler
func (a *Article) DeleteArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err = a.UcArticle.Delete(r.Context(), id)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	res.Message = "success delete article"

	return res, nil
}

func (a *Article) UpdateArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	req := RequestUpdateArticle{}
	b, err := rest.Bind(r, &req)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}
	if err = b.Validate(); err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	id, _ := strconv.Atoi(req.ID)

	// get user
	u, err := a.UcUser.GetUserID(r.Context(), req.UserID)
	if err != nil {
		return res, rest.ErrNotFound(w, r, err)
	}

	tags := []entity.Tag{}
	for _, tag := range req.Tags {
		t := entity.Tag{
			Name: tag,
		}

		tags = append(tags, t)
	}

	entity := entity.Article{
		ID:    id,
		Title: req.Title,
		Body:  req.Body,
		User:  &u,
		Tags:  tags,
	}

	article, err := a.UcArticle.Update(r.Context(), id, entity)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping request to entity
	res.Message = "success update article"
	res.Article = &article

	return res, nil
}

func (a *Article) GetByIDArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	e, err := a.UcArticle.GetByID(r.Context(), id)
	if err != nil {
		return res, rest.ErrNotFound(w, r, err)
	}

	// mapping entity to response
	res.Message = "success get article"
	res.Article = &e

	return res, nil
}

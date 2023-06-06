package article

import (
	"context"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent/article"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent/user"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"go.opentelemetry.io/otel/trace"
)

// GetAll returns resource users.
func (i *impl) GetAll(ctx context.Context, request entity.RequestGetArticles) (entity.ResponseGetArticles, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	client := i.adapter.PersistYmirBlog
	query := client.Article.Query()

	if request.Title != nil {
		query = query.Where(article.TitleContains(*request.Title))
	}

	if request.UserID != nil {
		query = query.Where(article.HasUserWith(user.IDEQ(*request.UserID)))
	}

	// pagination
	total, err := query.Count(ctx)
	if err != nil {
		return entity.ResponseGetArticles{}, err
	}
	metadata := rest.Pagination{
		Page:  request.Page,
		Limit: request.Limit,
		Total: total,
	}

	offset := (request.Page - 1) * request.Limit
	items, err := query.
		Limit(request.Limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return entity.ResponseGetArticles{}, err
	}

	res := entity.ResponseGetArticles{}
	for _, a := range items {
		entityArticle := entity.Article{
			Title: a.Title,
			Body:  a.Body,
		}
		res.Items = append(res.Items, &entityArticle)
	}
	res.Metadata = metadata

	return res, nil
}

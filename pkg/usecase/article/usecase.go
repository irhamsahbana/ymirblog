package article

import (
	"context"
	"net/http"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"go.opentelemetry.io/otel/trace"
)

// GetAll returns resource pokemon api.
func (i *impl) GetAll(ctx context.Context, r *http.Request, request entity.RequestGetArticles) ([]*entity.Article, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	// l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog
	query := client.Article.Query()

	// pagination
	total, err := query.Count(ctx)
	if err != nil {
		return []*entity.Article{}, err
	}

	rest.Paging(r, rest.Pagination{
		Page:  request.Page,
		Limit: request.Limit,
		Total: total,
	})

	offset := (request.Page - 1) * request.Limit
	items, err := query.
		Limit(request.Limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return []*entity.Article{}, err
	}

	res := []*entity.Article{}
	for _, a := range items {
		entityArticle := entity.Article{
			Title: a.Title,
			Body:  a.Body,
		}
		res = append(res, &entityArticle)
	}

	return res, nil
}

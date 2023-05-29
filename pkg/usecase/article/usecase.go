package article

import (
	"context"

	"github.com/rs/zerolog/log"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
	"go.opentelemetry.io/otel/trace"
)

// GetAll returns resource pokemon api.
func (i *impl) GetAll(ctx context.Context) ([]*entity.Article, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog
	articles, err := client.Article.Query().All(ctx)
	if err != nil {
		l.Error().Err(err).Msg("GetAll")
		return nil, err
	}

	res := []*entity.Article{}
	for _, a := range articles {
		entityArticle := entity.Article{
			Title: a.Title,
			Body:  a.Body,
		}
		res = append(res, &entityArticle)
	}

	return res, nil
}

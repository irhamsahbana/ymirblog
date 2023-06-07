package article

import (
	"context"

	"github.com/rs/zerolog/log"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent/article"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent/tag"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent/user"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
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

// GetByID returns resource article api.
func (i *impl) GetByID(ctx context.Context, id int) (e entity.Article, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog
	a, err := client.Article.Query().
		WithTags().
		WithUser().
		Where(
			article.ID(id),
		).
		First(ctx)
	if err != nil {
		l.Error().Err(err).Msg("GetByID")
		return e, err
	}

	e.ID = a.ID
	e.Title = a.Title
	e.Body = a.Body
	e.User = &entity.User{
		ID:    a.Edges.User.ID,
		Name:  a.Edges.User.Name,
		Email: a.Edges.User.Email,
	}

	for _, t := range a.Edges.Tags {
		e.Tags = append(e.Tags, entity.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return e, nil
}

// Create returns resource article api.
func (i *impl) Create(ctx context.Context, e entity.Article) (entity.Article, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog
	res := &e

	if err := client.WithTransaction(ctx, func(ctx context.Context, tx *ent.Tx) error {
		// create entity article, with user id
		eDB, err := tx.Article.Create().
			SetTitle(e.Title).
			SetBody(e.Body).
			SetUserID(e.User.ID).
			Save(ctx)
		if err != nil {
			l.Error().Err(err).Msg("Create")
			return err
		}

		// create entity tag
		for _, t := range e.Tags {
			findTag, err := tx.Tag.Query().Where(
				tag.Name(t.Name),
			).First(ctx)

			if err != nil {
				if ent.IsNotFound(err) { // if tag not found, create new tag and add to article
					tagDB, err := tx.Tag.Create().
						SetName(t.Name).
						Save(ctx)
					if err != nil {
						l.Error().Err(err).Msg("Create")
						return err
					}

					_, err = tx.Article.UpdateOneID(eDB.ID).
						AddTags(tagDB).
						Save(ctx)
					if err != nil {
						l.Error().Err(err).Msg("Create")
						return err
					}
				} else { // if error is not found error, rollback
					l.Error().Err(err).Msg("Create")
					return err
				}
			} else { // if tag found in DB, add tag to article
				_, err = tx.Article.UpdateOneID(eDB.ID).
					AddTags(findTag).
					Save(ctx)
				if err != nil {
					l.Error().Err(err).Msg("Create")
					return err
				}
			}
		}

		// find article with tags
		eDB, err = tx.Article.Query().
			Where(
				article.ID(eDB.ID),
			).
			WithTags().
			WithUser().
			First(ctx)
		if err != nil {
			l.Error().Err(err).Msg("Create")
			return err
		}

		// convert to entity
		res.ID = eDB.ID
		res.Title = eDB.Title
		res.Body = eDB.Body
		res.Tags = []entity.Tag{}
		for _, t := range eDB.Edges.Tags {
			res.Tags = append(res.Tags, entity.Tag{
				ID:   t.ID,
				Name: t.Name,
			})
		}
		res.User = &entity.User{
			ID:    eDB.Edges.User.ID,
			Name:  eDB.Edges.User.Name,
			Email: eDB.Edges.User.Email,
		}

		return nil
	}); err != nil {
		return e, err
	}

	return *res, nil
}

// Delete returns resource article api.
func (i *impl) Delete(ctx context.Context, id int) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog

	// delete article with tags
	err := client.Article.DeleteOneID(id).Exec(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Delete")
		return err
	}

	return nil
}

func (i *impl) Update(ctx context.Context, id int, e entity.Article) (entity.Article, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	client := i.adapter.PersistYmirBlog
	res := &entity.Article{}

	// start transaction
	if err := client.WithTransaction(ctx, func(ctx context.Context, tx *ent.Tx) error {
		a, err := tx.Article.UpdateOneID(id).
			SetUserID(e.User.ID).
			SetTitle(e.Title).
			SetBody(e.Body).
			ClearTags().
			Save(ctx)
		if err != nil {
			l.Error().Err(err).Msg("Update")
			return err
		}

		// create entity tag
		for _, t := range e.Tags {
			findTag, err := tx.Tag.Query().Where(
				tag.Name(t.Name),
			).First(ctx)

			if err != nil {
				if ent.IsNotFound(err) { // if tag not found, create new tag and add to article
					_, err := tx.Tag.Create().
						SetName(t.Name).
						AddArticles(a).
						Save(ctx)
					if err != nil {
						l.Error().Err(err).Msg("Update")
						return err
					}
				} else { // if error is not found error, rollback
					l.Error().Err(err).Msg("Update")
					return err
				}
			} else { // if tag found in DB, add tag to article
				_, err = tx.Article.UpdateOneID(a.ID).
					AddTags(findTag).
					Save(ctx)
				if err != nil {
					l.Error().Err(err).Msg("Update")
					return err
				}
			}
		}

		// find article with tags
		aDB, err := tx.Article.Query().
			Where(
				article.ID(id),
			).
			WithTags().
			WithUser().
			First(ctx)
		if err != nil {
			l.Error().Err(err).Msg("Update")
			return err
		}

		// convert to entity
		res.ID = aDB.ID
		res.Title = aDB.Title
		res.Body = aDB.Body

		res.Tags = []entity.Tag{}
		for _, t := range aDB.Edges.Tags {
			res.Tags = append(res.Tags, entity.Tag{
				ID:   t.ID,
				Name: t.Name,
			})
		}

		res.User = &entity.User{
			ID:    aDB.Edges.User.ID,
			Name:  aDB.Edges.User.Name,
			Email: aDB.Edges.User.Email,
		}

		return nil
	}); err != nil {
		return e, err
	}

	return *res, nil
}

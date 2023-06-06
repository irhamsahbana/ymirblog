package user

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
	"go.opentelemetry.io/otel/trace"
)

func (i *impl) CreateUser(ctx context.Context, newUser entity.User) (entity.User, error) {

	// validate persistence connection
	if i.adapter.PersistYmirBlog == nil {
		return newUser, errors.New("ymir blog persistence connection is nil")
	}

	//create user
	entUser, err := i.adapter.PersistYmirBlog.User.Create().
		SetName(newUser.Name).
		SetEmail(newUser.Email).
		Save(ctx)

	// mapping *ent.User to entity.User
	newUser.ID = entUser.ID

	return newUser, err
}

// Get User By Id
func (i *impl) GetUserID(ctx context.Context, ID int) (entity.User, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	user, err := i.adapter.PersistYmirBlog.User.Get(ctx, ID)
	if err != nil {
		l.Error().Err(err).Msg("GetBy ID")
		return entity.User{}, err
	}

	userID := entity.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return userID, nil

}

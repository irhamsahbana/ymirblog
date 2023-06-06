package user

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
	"go.opentelemetry.io/otel/trace"
)

// create User Usecase
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

// Update User Use case
func (i *impl) UpdateUser(ctx context.Context, ID int, updateUser entity.User) (entity.User, error) {
	// Update User
	entUser, err := i.adapter.PersistYmirBlog.User.UpdateOneID(ID).
		SetName(updateUser.Name).
		SetEmail(updateUser.Email).
		Save(ctx)
	if err != nil {
		return entity.User{}, err
	}

	updateUser = entity.User{
		ID:    entUser.ID,
		Name:  entUser.Name,
		Email: entUser.Email,
	}

	return updateUser, err
}

// get all user usecase
func (i *impl) GetAllUser(ctx context.Context) ([]entity.User, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	users, err := i.adapter.PersistYmirBlog.User.Query().All(ctx)
	if err != nil {
		l.Error().Err(err).Msg("GetAll")
		return nil, err
	}

	getAllUser := []entity.User{}
	for _, user := range users {
		entityUser := entity.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		getAllUser = append(getAllUser, entityUser)
	}

	return getAllUser, nil
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

// Delete User
func (i *impl) DeleteUser(ctx context.Context, ID int) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	// validate persistence connection
	if i.adapter.PersistYmirBlog == nil {
		return errors.New("ymir blog persistence connection is nil")
	}

	err := i.adapter.PersistYmirBlog.User.DeleteOneID(ID).Exec(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Delete ID")
		return err
	}

	return err
}

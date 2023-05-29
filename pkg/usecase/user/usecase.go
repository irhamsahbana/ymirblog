package user

import (
	"context"
	"errors"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
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

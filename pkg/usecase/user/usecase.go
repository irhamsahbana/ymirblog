package user

import (
	"context"
	"errors"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog"
)

func (i *impl) CreateUser(ctx context.Context, DB *ymirblog.Database, newUser entity.User) (entity.User, error) {

	// validate db connection
	if DB == nil {
		return newUser, errors.New("db connection is nil")
	}

	//create user
	entUser, err := DB.User.Create().
		SetName(newUser.Name).
		SetEmail(newUser.Email).
		Save(ctx)

	// mapping *ent.User to entity.User
	newUser.ID = entUser.ID

	return newUser, err
}

package user

import (
	"context"
	"errors"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
)

func (i *impl) CreateUser(ctx context.Context, newUser entity.User) (entity.User, error) {

	// validate db connection
	if i.adapter.PesistYmirBlog == nil {
		return newUser, errors.New("db connection is nil")
	}

	//create user
	entUser, err := i.adapter.PesistYmirBlog.User.Create().
		SetName(newUser.Name).
		SetEmail(newUser.Email).
		Save(ctx)

	// mapping *ent.User to entity.User
	newUser.ID = entUser.ID

	return newUser, err
}

package user

import (
	"context"
	"fmt"

	"reflect"

	"entgo.io/ent/dialect"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/adapters"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "user",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all user Component implementations.
type T interface {
	CreateUser(ctx context.Context, newUser entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, ID int, updateUser entity.User) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	GetUserID(ctx context.Context, ID int) (entity.User, error)
	DeleteUser(ctx context.Context, ID int) error
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a user Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}

func WithYmirBlogPersist() adapters.Option {
	return func(a *adapters.Adapter) {
		// adapter ymirblog sqlite
		if a.YmirBlogMySQL == nil {
			panic(fmt.Errorf("%s is not found", "YmirBlogMySQL"))
		}
		// persist ymirblog driver
		var c = ymirblog.Driver(
			ymirblog.WithDriver(a.YmirBlogMySQL, dialect.MySQL),
		)

		a.PersistYmirBlog = c
	}
}

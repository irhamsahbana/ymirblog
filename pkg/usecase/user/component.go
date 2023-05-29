package user

import (
	"context"

	"reflect"

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
	CreateUser(ctx context.Context, DB *ymirblog.Database, newUser entity.User) (entity.User, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a user Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}

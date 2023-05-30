// Package article is implements component logic.
package article

import (
	"context"
	"net/http"
	"reflect"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/adapters"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "article",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all article Component implementations.
type T interface {
	GetAll(ctx context.Context, r *http.Request, request entity.RequestGetArticles) ([]*entity.Article, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a article Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}

// Package rest is port handler.
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase/user"
)

// User handler instance data.
type User struct {
	UserUsecase user.T
	DB          *ymirblog.Database
}

// Register is endpoint group for handler.
func (u *User) Register(router chi.Router) {
	// PLEASE EDIT THIS EXAMPLE, how to register handler to router
	router.Get("/hello", rest.HandlerAdapter(u.User).JSON)
	router.Get("/hello-csv", rest.HandlerAdapter(u.UserCSV).CSV)

	// create handler for create user POST
	router.Post("/user/create", rest.HandlerAdapter(u.CreateUser).JSON)
}

// ResponseUser User handler response. /** PLEASE EDIT THIS EXAMPLE, return handler response */.
type ResponseUser struct {
	Message string
	User    entity.User
}

// Create User handler
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {

	request := entity.User{}
	b, err := rest.Bind[entity.User](r, &request)
	if err != nil {
		return ResponseUser{}, rest.ErrBadRequest(w, r, err)
	}
	if err = b.Validate(); err != nil {
		return ResponseUser{}, rest.ErrBadRequest(w, r, err)
	}

	//create user
	user, err := u.UserUsecase.CreateUser(r.Context(), request)
	if err != nil {
		return ResponseUser{
			Message: err.Error(),
		}, rest.ErrBadRequest(w, r, err)
	}

	return ResponseUser{
		Message: "success",
		User:    user,
	}, nil

}

// User endpoint func. /** PLEASE EDIT THIS EXAMPLE, return handler response */.
func (u *User) User(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {
	_, span, l := tracer.StartSpanLogTrace(r.Context(), "User")
	defer span.End()

	l.Info().Str("Hello", "World").Msg("this")

	return ResponseUser{
		Message: "Hello everybody",
	}, nil
}

// UserCSV endpoint func. /** PLEASE EDIT THIS EXAMPLE, return handler response */.
func (u *User) UserCSV(w http.ResponseWriter, r *http.Request) (rest.ResponseCSV, error) {
	_, span, l := tracer.StartSpanLogTrace(r.Context(), "UserCSV")
	defer span.End()

	l.Info().Str("Hello", "World").Msg("this")

	rows := make([][]string, 0)
	rows = append(rows, []string{"SO Number", "Nama Warung", "Area", "Fleet Number", "Jarak Warehouse", "Urutan"})
	rows = append(rows, []string{"SO45678", "WPD00011", "Jakarta Selatan", "1", "45.00", "1"})
	rows = append(rows, []string{"SO45645", "WPD001123", "Jakarta Selatan", "1", "43.00", "2"})
	rows = append(rows, []string{"SO45645", "WPD003343", "Jakarta Selatan", "1", "43.00", "3"})
	return rest.ResponseCSV{
		Filename: "warehouse",
		Rows:     rows,
	}, nil
}

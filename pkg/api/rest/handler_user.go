// Package rest is port handler.
package rest

import (
	"encoding/json"
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
	router.Post("/user/create", u.CreateUser)
}

// ResponseUser User handler response. /** PLEASE EDIT THIS EXAMPLE, return handler response */.
type ResponseUser struct {
	Message string
}

// Create User handler
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {

	var newUser entity.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//validation user name
	if newUser.Name == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//validation user email
	if newUser.Email == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//create user
	user, err := u.UserUsecase.CreateUser(r.Context(), u.DB, newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Write the user JSON data to the response
	w.Write(userJSON)

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

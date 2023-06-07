// Package rest is port handler.
package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/usecase/user"
)

// User handler instance data.
type User struct {
	UserUsecase user.T
}

// Register is endpoint group for handler.
func (u *User) Register(router chi.Router) {
	// create handler for user
	router.Post("/users", rest.HandlerAdapter(u.CreateUser).JSON)
	router.Patch("/users/{id}", rest.HandlerAdapter(u.UpdateUser).JSON)
	router.Get("/users", rest.HandlerAdapter(u.GetAllUser).JSON)
	router.Get("/users/{id}", rest.HandlerAdapter(u.GetUserID).JSON)
	router.Delete("/users/{id}", rest.HandlerAdapter(u.DeleteUser).JSON)
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
		User:    &user,
	}, nil
}

// Update User Handler
func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {
	//update user
	ID := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(ID)

	request := RequestUser{}
	b, err := rest.Bind(r, &request)
	fmt.Println(err)
	if err != nil {
		return ResponseUser{}, rest.ErrBadRequest(w, r, err)
	}
	if err = b.Validate(); err != nil {
		return ResponseUser{}, rest.ErrBadRequest(w, r, err)
	}

	entityUpdateUser := entity.User{
		Name:  request.Name,
		Email: request.Email,
	}

	userUpdate, err := u.UserUsecase.UpdateUser(r.Context(), id, entityUpdateUser)
	if err != nil {
		return ResponseUser{
			Message: err.Error(),
		}, rest.ErrBadRequest(w, r, err)
	}

	return ResponseUser{
		Message: "success",
		User:    &userUpdate,
	}, nil
}

// GetAllArticle Handler
func (u *User) GetAllUser(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {
	users, err := u.UserUsecase.GetAllUser(r.Context())
	if err != nil {
		return ResponseUser{}, err
	}

	return ResponseUser{
		Message: "succes",
		Users:   users,
	}, nil
}

func (u *User) GetUserID(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {
	ID := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(ID)

	userID, err := u.UserUsecase.GetUserID(r.Context(), id)
	if err != nil {
		return ResponseUser{
			Message: err.Error(),
		}, rest.ErrBadRequest(w, r, err)
	}

	return ResponseUser{
		Message: "success",
		User:    &userID,
	}, nil
}

func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request) (ResponseUser, error) {
	ID := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(ID)

	err := u.UserUsecase.DeleteUser(r.Context(), id)
	if err != nil {
		return ResponseUser{
			Message: err.Error(),
		}, rest.ErrBadRequest(w, r, err)
	}

	return ResponseUser{
		Message: "success",
	}, nil
}

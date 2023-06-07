package rest

import "gitlab.playcourt.id/dedenurr12/ymirblog/pkg/entity"

type RequestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email,required"`
}

type ResponseUser struct {
	Message string        `json:"message"`
	User    *entity.User  `json:"user,omitempty"`
	Users   []entity.User `json:"users,omitempty"`
}

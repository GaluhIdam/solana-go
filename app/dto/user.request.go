package dto

import (
	"github.com/go-playground/validator/v10"

)

type UserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=3,max=100"`
}

func (u *UserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

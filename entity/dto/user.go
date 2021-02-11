package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserSignUpRequest struct {
	Email         string `json:"email" validate:"required,email"`
	PlainPassword string `json:"password" validate:"required,min=5"`
}

func (u *UserSignUpRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type UserSignUpResponse struct {
	ID        string    `json:"id" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

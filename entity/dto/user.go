package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// TODO, add validate tag
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
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

type UserLoginRequest struct {
	Email         string `json:"email" validate:"required"`
	PlainPassword string `json:"password" validate:"required"`
}

func (u *UserLoginRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type UserRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (u *UserRefreshRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

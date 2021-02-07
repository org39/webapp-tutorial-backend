package dto

import (
	"github.com/go-playground/validator/v10"
)

type AuthTokenPair struct {
	AccessToken   string `json:"token"`
	RefereshToken string `json:"referesh_token"`
}

type AuthGenerateRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (u *AuthGenerateRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type AuthRefereshRequest struct {
	RefereshToken string `json:"referesh_token" validate:"required"`
}

func (u *AuthRefereshRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type AuthVerifyRequest struct {
	AccessToken string `json:"token" validate:"required"`
}

func (u *AuthVerifyRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

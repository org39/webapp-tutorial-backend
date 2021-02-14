package dto

import (
	"github.com/go-playground/validator/v10"
)

type AuthTokenPair struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthGenerateRequest struct {
	ID string `json:"id" validate:"required,uuid4"`
}

func (u *AuthGenerateRequest) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (u *AuthRefreshRequest) Valid() error {
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

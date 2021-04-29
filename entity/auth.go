package entity

import (
	"github.com/go-playground/validator/v10"
)

type AuthTokenPair struct {
	AccessToken  string `validate:"required"`
	RefreshToken string `validate:"required"`
}

func (u *AuthTokenPair) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

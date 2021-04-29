package entity

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (f *Validator) ValidateEmail(email string) error {
	err := validator.New().Var(email, "required,email")
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (f *Validator) ValidatePlainPassword(password string) error {
	err := validator.New().Var(password, "required,min=5")
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (f *Validator) ValidateID(id string) error {
	err := validator.New().Var(id, "required,uuid4")
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (f *Validator) ValidateToken(token string) error {
	err := validator.New().Var(token, "required")
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

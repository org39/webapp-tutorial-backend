package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        string    `validate:"required,uuid4"`
	Email     string    `validate:"required,email"`
	Password  string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
}

func (u *User) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

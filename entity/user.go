package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/org39/webapp-tutorial-backend/pkg/crypt"
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

func (u *User) ValidPassword(plainPassword string) error {
	return crypt.Compare(u.Password, []byte(plainPassword))
}

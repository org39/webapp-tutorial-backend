package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Todo struct {
	ID        string `validate:"required,uuid4"`
	UserID    string `validate:"required,uuid4"`
	Content   string `validate:"required"`
	Completed bool
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	Deleted   bool
}

func (u *Todo) Valid() error {
	err := validator.New().Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

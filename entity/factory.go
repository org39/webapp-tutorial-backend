package entity

import (
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"

	"github.com/org39/webapp-tutorial-backend/pkg/crypt"
	"github.com/org39/webapp-tutorial-backend/pkg/uuid"
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewUser(email string, plainPassword string, createdAt time.Time) (*User, error) {
	uuid, err := uuid.New()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := crypt.Hash([]byte(plainPassword))
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: createdAt,
	}, nil
}

func (f *Factory) FromUserDTO(u dto.User) (*User, error) {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}, nil
}

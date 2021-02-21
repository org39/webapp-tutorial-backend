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

func (f *Factory) NewUser(email string, plainPassword string) (*User, error) {
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
		CreatedAt: time.Now(),
	}, nil
}

func (f *Factory) FromUserDTO(u *dto.User) (*User, error) {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (f *Factory) NewTodo(user *User, content string) (*Todo, error) {
	uuid, err := uuid.New()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &Todo{
		ID:        uuid,
		UserID:    user.ID,
		Content:   content,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
		Deleted:   false,
	}, nil
}

func (f *Factory) FromTodoDTO(d *dto.Todo) (*Todo, error) {
	return &Todo{
		ID:        d.ID,
		UserID:    d.UserID,
		Content:   d.Content,
		Completed: d.Completed,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Deleted:   d.Deleted,
	}, nil
}

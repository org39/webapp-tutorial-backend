package dto

import "time"

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewUser(id string, email string, password string, createdAt time.Time) *User {
	return &User{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	}
}

func (f *Factory) NewTodo(id string, userID string, content string, completed bool, createdAt time.Time, updatedAt time.Time, deleted bool) *Todo {
	return &Todo{
		ID:        id,
		UserID:    userID,
		Content:   content,
		Completed: completed,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Deleted:   deleted,
	}
}

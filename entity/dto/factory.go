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

func (f *Factory) NewUserSignUpResponse(id string, email string, createdAt time.Time) *UserSignUpResponse {
	return &UserSignUpResponse{
		ID:        id,
		Email:     email,
		CreatedAt: createdAt,
	}
}

func (f *Factory) NewUserSignUpRequest(email string, plainPassword string) *UserSignUpRequest {
	return &UserSignUpRequest{
		Email:         email,
		PlainPassword: plainPassword,
	}
}

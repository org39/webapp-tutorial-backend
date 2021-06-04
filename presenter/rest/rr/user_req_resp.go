package rr

import (
	"time"
)

func (f *Factory) NewUserSignUpRequest(email string, plainPassword string) *UserSignUpRequest {
	return &UserSignUpRequest{
		Email:         email,
		PlainPassword: plainPassword,
	}
}

func (f *Factory) NewUserSignUpResponse(email string, createdAt time.Time, accessToken string) *UserSignUpResponse {
	return &UserSignUpResponse{
		Email:       email,
		AccessToken: accessToken,
		CreatedAt:   createdAt,
	}
}

func (f *Factory) NewUserLoginRequest(email string, plainPassword string) *UserLoginRequest {
	return &UserLoginRequest{
		Email:         email,
		PlainPassword: plainPassword,
	}
}

func (f *Factory) NewUserLoginResponse(accessToken string) *UserLoginResponse {
	return &UserLoginResponse{
		AccessToken: accessToken,
	}
}

func (f *Factory) NewUserRefreshRequest(refreshToken string) *UserRefreshRequest {
	return &UserRefreshRequest{
		RefreshToken: refreshToken,
	}
}

func (f *Factory) NewUserRefreshResponse(accessToken string) *UserRefreshResponse {
	return &UserRefreshResponse{
		AccessToken: accessToken,
	}
}

func (f *Factory) NewUserResponse(email string, createdAt time.Time) *UserResponse {
	return &UserResponse{
		Email:     email,
		CreatedAt: createdAt,
	}
}

// ------------------------------------------------------------------
type UserSignUpRequest struct {
	Email         string `json:"email"`
	PlainPassword string `json:"password"`
}

type UserSignUpResponse struct {
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserLoginRequest struct {
	Email         string `json:"email"`
	PlainPassword string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserRefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

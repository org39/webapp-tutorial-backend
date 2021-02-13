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

func (f *Factory) NewAuthGenerateRequest(email string) *AuthGenerateRequest {
	return &AuthGenerateRequest{
		Email: email,
	}
}

func (f *Factory) NewAuthTokenPair(token string, refreshToken string) *AuthTokenPair {
	return &AuthTokenPair{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
}

func (f *Factory) NewAuthRefreshRequest(refreshToken string) *AuthRefreshRequest {
	return &AuthRefreshRequest{
		RefreshToken: refreshToken,
	}
}

func (f *Factory) NewAuthVerifyRequest(accessToken string) *AuthVerifyRequest {
	return &AuthVerifyRequest{
		AccessToken: accessToken,
	}
}

func (f *Factory) NewUserLoginRequest(email string, plainPassword string) *UserLoginRequest {
	return &UserLoginRequest{
		Email:         email,
		PlainPassword: plainPassword,
	}
}

func (f *Factory) NewUserRefreshRequest(refreshToken string) *UserRefreshRequest {
	return &UserRefreshRequest{
		RefreshToken: refreshToken,
	}
}

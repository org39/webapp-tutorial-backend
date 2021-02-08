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

func (f *Factory) NewUserSignUpResponse(id string, email string, accessToken string, refereshToken string, createdAt time.Time) *UserSignUpResponse {
	return &UserSignUpResponse{
		ID:            id,
		Email:         email,
		AccessToken:   accessToken,
		RefereshToken: refereshToken,
		CreatedAt:     createdAt,
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

func (f *Factory) NewAuthTokenPair(token string, refereshToken string) *AuthTokenPair {
	return &AuthTokenPair{
		AccessToken:   token,
		RefereshToken: refereshToken,
	}
}

func (f *Factory) NewAuthRefereshRequest(refereshToken string) *AuthRefereshRequest {
	return &AuthRefereshRequest{
		RefereshToken: refereshToken,
	}
}

func (f *Factory) NewAuthVerifyRequest(accessToken string) *AuthVerifyRequest {
	return &AuthVerifyRequest{
		AccessToken: accessToken,
	}
}

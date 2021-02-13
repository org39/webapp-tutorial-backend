package user

//go:generate mockery --all

import (
	"context"
	"errors"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrSystemError    = errors.New("system error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrDatabaseError  = errors.New("database error")
)

type Usecase interface {
	SignUp(ctx context.Context, req *dto.UserSignUpRequest) (*dto.UserSignUpResponse, *dto.AuthTokenPair, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.AuthTokenPair, error)
	Refresh(ctx context.Context, req *dto.UserRefreshRequest) (*dto.AuthTokenPair, error)
}

type Repository interface {
	FetchByID(ctx context.Context, id string) (*dto.User, error)
	FetchByEmail(ctx context.Context, email string) (*dto.User, error)
	Store(ctx context.Context, u *dto.User) error
	Update(ctx context.Context, u *dto.User) error
}

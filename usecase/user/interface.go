package user

//go:generate mockery --all

import (
	"context"
	"errors"

	"github.com/org39/webapp-tutorial-backend/entity"
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
	FetchByID(ctx context.Context, id string) (*entity.User, error)
	SignUp(ctx context.Context, email string, plainPassword string) (*entity.User, *entity.AuthTokenPair, error)
	Login(ctx context.Context, email string, password string) (*entity.AuthTokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*entity.AuthTokenPair, error)
}

type Repository interface {
	FetchByID(ctx context.Context, id string) (*dto.User, error)
	FetchByEmail(ctx context.Context, email string) (*dto.User, error)
	Store(ctx context.Context, u *dto.User) error
	Update(ctx context.Context, u *dto.User) error
}

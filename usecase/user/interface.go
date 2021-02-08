package user

//go:generate mockery --all

import (
	"context"
	"errors"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
)

var (
	ErrInvalidSignUpReq = errors.New("invalid signup request")
	ErrNotFound         = errors.New("not found")
	ErrSystemError      = errors.New("system error")
	ErrDatabaseError    = errors.New("database error")
)

type Usecase interface {
	SignUp(ctx context.Context, req *dto.UserSignUpRequest) (*dto.UserSignUpResponse, error)
}

type Repository interface {
	FetchByID(ctx context.Context, id string) (*dto.User, error)
	FetchByEmail(ctx context.Context, email string) (*dto.User, error)
	Store(ctx context.Context, u *dto.User) error
	Update(ctx context.Context, u *dto.User) error
}

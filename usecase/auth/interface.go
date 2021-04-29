package auth

//go:generate mockery --all

import (
	"context"
	"errors"

	"github.com/org39/webapp-tutorial-backend/entity"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrSystemError    = errors.New("system error")
	ErrUnauthorized   = errors.New("unauthorized")
)

type Usecase interface {
	GenereateToken(ctx context.Context, id string) (*entity.AuthTokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.AuthTokenPair, error)
	VerifyToken(ctx context.Context, accessToken string) (string, error)
}

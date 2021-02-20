package auth

//go:generate mockery --all

import (
	"context"
	"errors"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrSystemError    = errors.New("system error")
	ErrUnauthorized   = errors.New("unauthorized")
)

type Usecase interface {
	GenereateToken(ctx context.Context, req *dto.AuthGenerateRequest) (*dto.AuthTokenPair, error)
	RefreshToken(ctx context.Context, req *dto.AuthRefreshRequest) (*dto.AuthTokenPair, error)
	VerifyToken(ctx context.Context, req *dto.AuthVerifyRequest) (string, error)
}

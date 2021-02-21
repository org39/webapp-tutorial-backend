package todo

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
	Create(ctx context.Context, u *dto.User, content string) (*dto.Todo, error)
	FetchAllByUser(ctx context.Context, u *dto.User) ([]*dto.Todo, error)
	FetchByID(ctx context.Context, u *dto.User, id string) (*dto.Todo, error)
	Update(ctx context.Context, u *dto.User, t *dto.Todo) (*dto.Todo, error)
	Delete(ctx context.Context, u *dto.User, t *dto.Todo) error
}

type Repository interface {
	Store(ctx context.Context, t *dto.Todo) error
	Update(ctx context.Context, t *dto.Todo) error
	Delete(ctx context.Context, t *dto.Todo) error
	FetchAllByUser(ctx context.Context, u *dto.User) ([]*dto.Todo, error)
	FetchByID(ctx context.Context, id string) (*dto.Todo, error)
}

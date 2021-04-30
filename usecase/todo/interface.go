package todo

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
	Create(ctx context.Context, user *entity.User, content string) (*entity.Todo, error)
	FetchAllByUser(ctx context.Context, user *entity.User, showCompleted bool, showDeleted bool) ([]*entity.Todo, error)
	FetchByID(ctx context.Context, user *entity.User, id string) (*entity.Todo, error)
	Update(ctx context.Context, user *entity.User, id string, content string, completed bool, deleted bool) (*entity.Todo, error)
	Delete(ctx context.Context, user *entity.User, id string) error
}

type Repository interface {
	Store(ctx context.Context, t *dto.Todo) error
	Update(ctx context.Context, t *dto.Todo) error
	Delete(ctx context.Context, t *dto.Todo) error
	FetchAllByUser(ctx context.Context, u *dto.User, showCompleted bool, showDeleted bool) ([]*dto.Todo, error)
	FetchByID(ctx context.Context, id string) (*dto.Todo, error)
}

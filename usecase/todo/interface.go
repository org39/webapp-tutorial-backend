package user

//go:generate mockery --all

import (
	// "context"
	"errors"
	// "github.com/org39/webapp-tutorial-backend/entity/dto"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrSystemError    = errors.New("system error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrDatabaseError  = errors.New("database error")
)

type Usecase interface {
}

type Repository interface {
}

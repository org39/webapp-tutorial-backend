package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity"
	"github.com/org39/webapp-tutorial-backend/entity/dto"
)

type Service struct {
	Repository Repository `inject:""`
}

func NewService(options ...func(*Service) error) (Usecase, error) {
	u := &Service{}

	for _, option := range options {
		if err := option(u); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func WithRepository(r Repository) func(*Service) error {
	return func(u *Service) error {
		u.Repository = r
		return nil
	}
}

func (u *Service) SignUp(ctx context.Context, req *dto.UserSignUpRequest) (*dto.UserSignUpResponse, error) {
	// test some validation on req
	if err := req.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid signup request: %w", err, ErrInvalidSignUpReq)
	}

	// test email alread exist
	_, err := u.Repository.FetchByEmail(ctx, req.Email)
	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("email already exist: %w", ErrInvalidSignUpReq)
	}

	// create user object
	user, err := entity.NewFactory().NewUser(req.Email, req.PlainPassword, time.Now())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrSystemError)
	}

	// validation user object
	if err := user.Valid(); err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrInvalidSignUpReq)
	}

	// store user
	userDTO := dto.NewFactory().NewUser(user.ID, user.Email, user.Password, user.CreatedAt)
	if err := u.Repository.Store(ctx, userDTO); err != nil {
		return nil, err
	}

	res := dto.NewFactory().NewUserSignUpResponse(user.ID, user.Email, user.CreatedAt)
	return res, nil
}

package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity"
	"github.com/org39/webapp-tutorial-backend/entity/dto"

	"github.com/org39/webapp-tutorial-backend/usecase/auth"
)

type Service struct {
	Repository  Repository   `inject:""`
	AuthUsecase auth.Usecase `inject:""`
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

func WithAuthUsecase(a auth.Usecase) func(*Service) error {
	return func(u *Service) error {
		u.AuthUsecase = a
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
	switch {
	case errors.Is(err, ErrNotFound):
		// do nothing
	case err == nil:
		return nil, fmt.Errorf("email already exist: %w", ErrInvalidSignUpReq)
	case err != nil:
		return nil, err
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

	token, err := u.AuthUsecase.GenereateToken(ctx, dto.NewFactory().NewAuthGenerateRequest(user.Email))
	if err != nil {
		return nil, toUserServiceError(err)
	}

	res := dto.NewFactory().NewUserSignUpResponse(user.ID, user.Email, token.AccessToken, token.RefereshToken, user.CreatedAt)
	return res, nil
}

func toUserServiceError(err error) error {
	switch {
	case errors.Is(err, auth.ErrInvalidRequest):
		return fmt.Errorf("%s: invalid signup request: %w", err, auth.ErrInvalidRequest)
	case errors.Is(err, auth.ErrSystemError):
		return fmt.Errorf("%s: %w", err, auth.ErrInvalidRequest)
	}

	return fmt.Errorf("%s: %w", err, auth.ErrInvalidRequest)
}

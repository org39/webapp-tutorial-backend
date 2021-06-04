package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/org39/webapp-tutorial-backend/entity"
	"github.com/org39/webapp-tutorial-backend/entity/dto"

	"github.com/org39/webapp-tutorial-backend/usecase/auth"
)

type Service struct {
	Repository   Repository   `inject:""`
	AuthUsecase  auth.Usecase `inject:""`
	PasswordSalt string       `inject:"usecase.user.password_salt"`
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

func (u *Service) SignUp(ctx context.Context, email string, plainPassword string) (*entity.User, *entity.AuthTokenPair, error) {
	// validation on parameters
	if err := entity.NewValidator().ValidateEmail(email); err != nil {
		return nil, nil, fmt.Errorf("%s: invalid signup request: %w", err, ErrInvalidRequest)
	}

	if err := entity.NewValidator().ValidatePlainPassword(plainPassword); err != nil {
		return nil, nil, fmt.Errorf("%s: invalid signup request: %w", err, ErrInvalidRequest)
	}

	// test email alread exist
	_, err := u.Repository.FetchByEmail(ctx, email)
	switch {
	case errors.Is(err, ErrNotFound):
		// do nothing
	case err == nil:
		return nil, nil, fmt.Errorf("email already exist: %w", ErrInvalidRequest)
	case err != nil:
		return nil, nil, err
	}

	// create user object
	saltedPassword := fmt.Sprintf("%s%s", plainPassword, u.PasswordSalt)
	user, err := entity.NewFactory().NewUser(email, saltedPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", err.Error(), ErrSystemError)
	}

	// validation user object
	if err := user.Valid(); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", err.Error(), ErrInvalidRequest)
	}

	// store user
	userDTO := dto.NewFactory().NewUser(user.ID, user.Email, user.Password, user.CreatedAt)
	if err := u.Repository.Store(ctx, userDTO); err != nil {
		return nil, nil, err
	}

	token, err := u.AuthUsecase.GenereateToken(ctx, user.ID)
	if err != nil {
		return nil, nil, toUserServiceError(err)
	}

	return user, token, nil
}

func (u *Service) Login(ctx context.Context, email string, plainPassword string) (*entity.AuthTokenPair, error) {
	// test email alread exist
	userDTO, err := u.Repository.FetchByEmail(ctx, email)
	switch {
	case errors.Is(err, ErrNotFound):
		return nil, fmt.Errorf("email not found: %w", ErrNotFound)
	case err != nil:
		return nil, err
	}

	user, err := entity.NewFactory().FromUserDTO(userDTO)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	saltedPassword := fmt.Sprintf("%s%s", plainPassword, u.PasswordSalt)
	if err := user.ValidPassword(saltedPassword); err != nil {
		return nil, fmt.Errorf("%w", ErrUnauthorized)
	}

	token, err := u.AuthUsecase.GenereateToken(ctx, user.ID)
	if err != nil {
		return nil, toUserServiceError(err)
	}

	return token, nil
}

func (u *Service) Refresh(ctx context.Context, refreshToken string) (*entity.AuthTokenPair, error) {
	token, err := u.AuthUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, toUserServiceError(err)
	}

	return token, nil
}

func (u *Service) FetchByID(ctx context.Context, id string) (*entity.User, error) {
	userDTO, err := u.Repository.FetchByID(ctx, id)
	if err != nil {
		return nil, toUserServiceError(err)
	}

	user, err := entity.NewFactory().FromUserDTO(userDTO)
	if err != nil {
		return nil, toUserServiceError(err)
	}

	return user, nil
}

func toUserServiceError(err error) error {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return fmt.Errorf("%s: %w", err, ErrUnauthorized)
	case errors.Is(err, auth.ErrInvalidRequest):
		return fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	case errors.Is(err, auth.ErrSystemError):
		return fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	return fmt.Errorf("%s: %w", err, ErrSystemError)
}

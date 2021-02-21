package todo

import (
	"context"
	"fmt"

	"github.com/org39/webapp-tutorial-backend/entity"
	"github.com/org39/webapp-tutorial-backend/entity/dto"
)

type Service struct {
	Repository Repository `inject:""`
}

func NewService(options ...func(*Service) error) (Usecase, error) {
	s := &Service{}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func WithRepository(r Repository) func(*Service) error {
	return func(s *Service) error {
		s.Repository = r
		return nil
	}
}

func (s *Service) Create(ctx context.Context, u *dto.User, content string) (*dto.Todo, error) {
	// test some validation on req
	if err := u.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	user, err := entity.NewFactory().FromUserDTO(u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	todo, err := entity.NewFactory().NewTodo(user, content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	// validation todo object
	if err := todo.Valid(); err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrInvalidRequest)
	}

	todoDTO := dto.NewFactory().NewTodo(todo.ID, todo.UserID, todo.Content, todo.Completed, todo.CreatedAt, todo.UpdatedAt, todo.Deleted)
	if err := s.Repository.Store(ctx, todoDTO); err != nil {
		return nil, err
	}

	return todoDTO, nil
}

func (s *Service) FetchAllByUser(ctx context.Context, u *dto.User) ([]*dto.Todo, error) {
	// test some validation on req
	if err := u.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	return s.Repository.FetchAllByUser(ctx, u)
}

func (s *Service) FetchByID(ctx context.Context, u *dto.User, id string) (*dto.Todo, error) {
	todo, err := s.Repository.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if u.ID != todo.UserID {
		return nil, ErrUnauthorized
	}

	return todo, nil
}

func (s *Service) Update(ctx context.Context, u *dto.User, t *dto.Todo) (*dto.Todo, error) {
	// test some validation on req
	if err := t.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	if u.ID != t.UserID {
		return nil, ErrUnauthorized
	}

	if err := s.Repository.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Delete(ctx context.Context, u *dto.User, t *dto.Todo) error {
	// test some validation on req
	if err := t.Valid(); err != nil {
		return fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	if u.ID != t.UserID {
		return ErrUnauthorized
	}

	// mark deleted
	t.Deleted = true
	return s.Repository.Update(ctx, t)
}

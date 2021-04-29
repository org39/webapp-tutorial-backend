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

func (s *Service) Create(ctx context.Context, user *entity.User, content string) (*entity.Todo, error) {
	// test some validation on req
	if err := user.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
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

	return todo, nil
}

func (s *Service) FetchAllByUser(ctx context.Context, user *entity.User) ([]*entity.Todo, error) {
	// test some validation on req
	if err := user.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	userDTO := dto.NewFactory().NewUser(user.ID, user.Email, user.Password, user.CreatedAt)
	todoDTOs, err := s.Repository.FetchAllByUser(ctx, userDTO)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrDatabaseError)
	}

	todos := make([]*entity.Todo, len(todoDTOs))
	for i, todoDTO := range todoDTOs {
		todo, err := entity.NewFactory().FromTodoDTO(todoDTO)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
		}

		todos[i] = todo
	}

	return todos, nil
}

func (s *Service) FetchByID(ctx context.Context, u *entity.User, id string) (*entity.Todo, error) {
	todoDTO, err := s.Repository.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo, err := entity.NewFactory().FromTodoDTO(todoDTO)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	if u.ID != todo.UserID {
		return nil, ErrUnauthorized
	}

	return todo, nil
}

func (s *Service) Update(ctx context.Context, user *entity.User, id string, content string, completed bool, deleted bool) (*entity.Todo, error) {
	// fetch todo
	ori, err := s.Repository.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.ID != ori.UserID {
		return nil, ErrUnauthorized
	}

	// create new Todo
	newTodo, err := entity.NewFactory().FromTodoDTO(ori)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrSystemError)
	}

	newTodo.Content = content
	newTodo.Completed = completed
	newTodo.Deleted = deleted

	// test some validation on new Todo
	if err := newTodo.Valid(); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", err, ErrInvalidRequest)
	}

	// Update
	newTodoDTO := dto.NewFactory().NewTodo(newTodo.ID, newTodo.UserID, newTodo.Content, newTodo.Completed, newTodo.CreatedAt, newTodo.UpdatedAt, newTodo.Deleted)
	if err := s.Repository.Update(ctx, newTodoDTO); err != nil {
		return nil, err
	}

	return newTodo, nil
}

func (s *Service) Delete(ctx context.Context, user *entity.User, id string) error {
	t, err := s.Repository.FetchByID(ctx, id)
	if err != nil {
		return err
	}

	if user.ID != t.UserID {
		return ErrUnauthorized
	}

	// mark deleted
	t.Deleted = true
	return s.Repository.Update(ctx, t)
}

package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"

	sq "github.com/Masterminds/squirrel"
)

var (
	todoCols = []string{"id", "user_id", "content", "completed", "created_at", "updated_at", "deleted"}
)

type TodoRepository struct {
	DB    *db.DB `inject:""`
	Table string `inject:"repo.todo.table"`
}

func NewTodoRepository(options ...func(*TodoRepository) error) (todo.Repository, error) {
	r := &TodoRepository{}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func WithTodoDB(db *db.DB) func(*TodoRepository) error {
	return func(r *TodoRepository) error {
		r.DB = db
		return nil
	}
}

func WithTodoTable(table string) func(*TodoRepository) error {
	return func(r *TodoRepository) error {
		r.Table = table
		return nil
	}
}

func (r *TodoRepository) Store(ctx context.Context, t *dto.Todo) error {
	query, args, err := sq.Insert(r.Table).Columns(todoCols...).
		Values(t.ID, t.UserID, t.Content, t.Completed, t.CreatedAt, t.UpdatedAt, t.Deleted).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}
	return nil
}

func (r *TodoRepository) Update(ctx context.Context, t *dto.Todo) error {
	query, args, err := sq.Update(r.Table).
		Set("content", t.Content).
		Set("completed", t.Completed).
		Set("deleted", t.Deleted).
		Where(sq.Eq{"id": t.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, t *dto.Todo) error {
	query, args, err := sq.Delete(r.Table).
		Where(sq.Eq{"id": t.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}
	return nil
}

func (r *TodoRepository) FetchAllByUser(ctx context.Context, u *dto.User, showCompleted bool, showDeleted bool) ([]*dto.Todo, error) {
	q := r.selectTodo()
	switch {
	case showCompleted && showDeleted:
		q = q.Where(sq.Eq{"user_id": u.ID})
	case !showCompleted && showDeleted:
		q = q.Where(sq.Eq{"user_id": u.ID, "completed": false})
	case !showCompleted && !showDeleted:
		q = q.Where(sq.Eq{"user_id": u.ID, "completed": false, "deleted": false})
	case showCompleted && !showDeleted:
		q = q.Where(sq.Eq{"user_id": u.ID, "deleted": false})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return []*dto.Todo{}, nil
	case err != nil:
		return nil, fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}
	defer rows.Close()

	todos := []*dto.Todo{}
	for rows.Next() {
		t, err := r.scanTodo(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (r *TodoRepository) FetchByID(ctx context.Context, id string) (*dto.Todo, error) {
	query, args, err := r.selectTodo().Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	row := r.DB.QueryRow(ctx, query, args...)
	t, err := r.scanTodo(row)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TodoRepository) selectTodo() sq.SelectBuilder {
	return sq.Select(todoCols...).From(r.Table)
}

func (r *TodoRepository) scanTodo(row db.Scanable) (*dto.Todo, error) {
	var id, userID, content string
	var completed, deleted bool
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &userID, &content, &completed, &createdAt, &updatedAt, &deleted)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, todo.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("%s: %w", err.Error(), todo.ErrDatabaseError)
	}

	return dto.NewFactory().NewTodo(id, userID, content, completed, createdAt, updatedAt, deleted), nil
}

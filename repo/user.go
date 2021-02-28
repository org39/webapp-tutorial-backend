package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	sq "github.com/Masterminds/squirrel"
)

var (
	userCols = []string{"id", "email", "password", "created_at"}
)

type UserRepository struct {
	DB    *db.DB `inject:""`
	Table string `inject:"repo.user.table"`
}

func NewUserRepository(options ...func(*UserRepository) error) (user.Repository, error) {
	r := &UserRepository{}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func WithUserDB(db *db.DB) func(*UserRepository) error {
	return func(r *UserRepository) error {
		r.DB = db
		return nil
	}
}

func WithUserTable(table string) func(*UserRepository) error {
	return func(r *UserRepository) error {
		r.Table = table
		return nil
	}
}

func (r *UserRepository) FetchByID(ctx context.Context, id string) (*dto.User, error) {
	query, args, err := r.selectUser().Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}

	row := r.DB.QueryRow(ctx, query, args...)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FetchByEmail(ctx context.Context, email string) (*dto.User, error) {
	query, args, err := r.selectUser().Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}

	row := r.DB.QueryRow(ctx, query, args...)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Store(ctx context.Context, u *dto.User) error {
	query, args, err := sq.Insert(r.Table).
		Columns(userCols...).
		Values(u.ID, u.Email, u.Password, u.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *dto.User) error {
	query, args, err := sq.Update(r.Table).
		Set("email", u.Email).
		Set("password", u.Password).
		Where(sq.Eq{"id": u.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}
	return nil
}

func (r *UserRepository) selectUser() sq.SelectBuilder {
	return sq.Select(userCols...).From(r.Table)
}

func (r *UserRepository) scanUser(row db.Scanable) (*dto.User, error) {
	var id, email, password string
	var CreatedAt time.Time

	err := row.Scan(&id, &email, &password, &CreatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, user.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("%s: %w", err.Error(), user.ErrDatabaseError)
	}

	return dto.NewFactory().NewUser(id, email, password, CreatedAt), nil
}

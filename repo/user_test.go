package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	UserRepository user.Repository
	DB             *db.DB
	Sqlmock        sqlmock.Sqlmock
}

func (s *UserRepoTestSuite) SetupTest() {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to sqlmock: %s", err))
	}
	s.DB = &db.DB{DB: mockdb}
	s.Sqlmock = mock

	r, err := NewUserRepository(
		WithTable("users"),
		WithDB(s.DB),
	)
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create repository: %s", err))
	}

	s.UserRepository = r
}

func (s *UserRepoTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *UserRepoTestSuite) TestFetchByEmailExist() {
	ctx := context.Background()
	email := "hatsune@miku.com"

	// mock database
	q := "SELECT id, email, password, created_at FROM users WHERE email = ?"
	s.Sqlmock.ExpectQuery(q).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).AddRow("id", email, "PASSWORD", time.Now()),
		)

	// assert
	user, err := s.UserRepository.FetchByEmail(ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), email, user.Email)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestFetchByEmailNotExist() {
	ctx := context.Background()
	email := "not-exist@mail.com"

	// mock database
	q := "SELECT id, email, password, created_at FROM users WHERE email = ?"
	s.Sqlmock.ExpectQuery(q).
		WithArgs(email).WillReturnError(sql.ErrNoRows)

	// assert
	u, err := s.UserRepository.FetchByEmail(ctx, email)
	assert.Nil(s.T(), u)
	assert.ErrorIs(s.T(), user.ErrNotFound, err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestStoreSuccess() {
	ctx := context.Background()
	u := dto.NewFactory().NewUser("5c2dd83a-6250-40f3-a47e-21d957c07d06", "hatsune@miku.com", "PASSWORD", time.Now())

	q := "INSERT INTO users"
	s.Sqlmock.ExpectBegin()
	s.Sqlmock.ExpectExec(q).
		WithArgs(u.ID, u.Email, u.Password, u.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Sqlmock.ExpectCommit()

	// assert
	err := s.UserRepository.Store(ctx, u)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestUpdateSuccess() {
	ctx := context.Background()
	u := dto.NewFactory().NewUser("5c2dd83a-6250-40f3-a47e-21d957c07d06", "hatsune@miku.com", "PASSWORD", time.Now())

	q := "UPDATE users"
	s.Sqlmock.ExpectBegin()
	s.Sqlmock.ExpectExec(q).
		WithArgs(u.Email, u.Password, u.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Sqlmock.ExpectCommit()

	// assert
	err := s.UserRepository.Update(ctx, u)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func TestUserRepo(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

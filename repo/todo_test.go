package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TodoRepoTestSuite struct {
	suite.Suite
	TodoRepository todo.Repository
	DB             *db.DB
	Sqlmock        sqlmock.Sqlmock
}

func (s *TodoRepoTestSuite) SetupTest() {
	mockdb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to sqlmock: %s", err))
	}
	s.DB = &db.DB{DB: mockdb}
	s.Sqlmock = mock

	r, err := NewTodoRepository(
		WithTodoTable("todos"),
		WithTodoDB(s.DB),
	)
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create repository: %s", err))
	}

	s.TodoRepository = r
}

func (s *TodoRepoTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *TodoRepoTestSuite) TestStoreSuccess() {
	ctx := context.Background()

	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	t := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	q := "INSERT INTO todos (id,user_id,content,completed,created_at,updated_at,deleted) VALUES (?,?,?,?,?,?,?)"
	s.Sqlmock.ExpectBegin()
	s.Sqlmock.ExpectExec(q).
		WithArgs(t.ID, t.UserID, t.Content, t.Completed, t.CreatedAt, t.UpdatedAt, t.Deleted).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Sqlmock.ExpectCommit()

	// assert
	err := s.TodoRepository.Store(ctx, t)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *TodoRepoTestSuite) TestUpdateSuccess() {
	ctx := context.Background()

	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	t := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	q := "UPDATE todos SET content = ?, completed = ?, deleted = ? WHERE id = ?"
	s.Sqlmock.ExpectBegin()
	s.Sqlmock.ExpectExec(q).
		WithArgs(t.Content, t.Completed, t.Deleted, t.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Sqlmock.ExpectCommit()

	// assert
	err := s.TodoRepository.Update(ctx, t)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *TodoRepoTestSuite) TestDeleteSuccess() {
	ctx := context.Background()

	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	t := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	q := "DELETE FROM todos WHERE id = ?"
	s.Sqlmock.ExpectBegin()
	s.Sqlmock.ExpectExec(q).
		WithArgs(t.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Sqlmock.ExpectCommit()

	// assert
	err := s.TodoRepository.Delete(ctx, t)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *TodoRepoTestSuite) TestFetchByIDExist() {
	ctx := context.Background()

	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	t := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	q := "SELECT id, user_id, content, completed, created_at, updated_at, deleted FROM todos WHERE id = ?"
	s.Sqlmock.ExpectQuery(q).
		WithArgs(t.ID).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "user_id", "content", "completed", "created_at", "updated_at", "deleted"}).
				AddRow(t.ID, t.UserID, t.Content, t.Completed, t.CreatedAt, t.UpdatedAt, t.Deleted),
		)

	// assert
	res, err := s.TodoRepository.FetchByID(ctx, t.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), t.ID, res.ID)
	assert.Equal(s.T(), t.UserID, res.UserID)
	assert.Equal(s.T(), t.Content, res.Content)
	assert.Equal(s.T(), t.Completed, res.Completed)
	assert.Equal(s.T(), t.Deleted, res.Deleted)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *TodoRepoTestSuite) TestFetchByIDNotExist() {
	ctx := context.Background()

	id := "4daaaea8-4721-4644-aaac-7958805b4530"

	q := "SELECT id, user_id, content, completed, created_at, updated_at, deleted FROM todos WHERE id = ?"
	s.Sqlmock.ExpectQuery(q).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	// assert
	res, err := s.TodoRepository.FetchByID(ctx, id)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), todo.ErrNotFound, err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func (s *TodoRepoTestSuite) TestFetchAllByUserNotExist() {
	ctx := context.Background()

	u := dto.NewFactory().NewUser("5c2dd83a-6250-40f3-a47e-21d957c07d06", "hatsune@miku.com", "PASSWORD", time.Now())
	q := "SELECT id, user_id, content, completed, created_at, updated_at, deleted FROM todos WHERE completed = ? AND deleted = ? AND user_id = ?"
	s.Sqlmock.ExpectQuery(q).
		WithArgs(false, false, u.ID).
		WillReturnError(sql.ErrNoRows)

	// assert
	res, err := s.TodoRepository.FetchAllByUser(ctx, u, false, false)
	assert.NotNil(s.T(), res)
	assert.Empty(s.T(), res)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.Sqlmock.ExpectationsWereMet())
}

func TestTodoRepo(t *testing.T) {
	suite.Run(t, new(TodoRepoTestSuite))
}

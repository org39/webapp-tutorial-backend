package todo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/todo/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TodoServiceTestSuite struct {
	suite.Suite
	Usecase    Usecase
	Repository *mocks.Repository
}

func (s *TodoServiceTestSuite) SetupTest() {
	s.Repository = new(mocks.Repository)

	usecase, err := NewService(
		WithRepository(s.Repository),
	)
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create usecase: %s", err))
	}

	s.Usecase = usecase
}

func (s *TodoServiceTestSuite) TestCreateSuccess() {
	ctx := context.Background()

	// mock repo
	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	userDTO := dto.NewFactory().NewUser(userID, "account@emai.com", "strong-password", time.Now())

	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	todoDTO := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	s.Repository.On("Store", ctx, mock.AnythingOfType("*dto.Todo")).Return(nil)

	// assert
	res, err := s.Usecase.Create(ctx, userDTO, todoDTO.Content)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), todoDTO.UserID, res.UserID)
}

func (s *TodoServiceTestSuite) TestFetctByIDSuccess() {
	ctx := context.Background()

	// mock repo
	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	userDTO := dto.NewFactory().NewUser(userID, "account@emai.com", "strong-password", time.Now())

	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	todoDTO := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	s.Repository.On("FetchByID", ctx, id).Return(todoDTO, nil)

	// assert
	res, err := s.Usecase.FetchByID(ctx, userDTO, id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), todoDTO.UserID, res.UserID)
}

func (s *TodoServiceTestSuite) TestFetctAllByUserSuccess() {
	ctx := context.Background()

	// mock repo
	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	userDTO := dto.NewFactory().NewUser(userID, "account@emai.com", "strong-password", time.Now())

	id0 := "4daaaea8-4721-4644-aaac-7958805b4530"
	todoDTO0 := dto.NewFactory().NewTodo(id0, userID, "things todo", false, time.Now(), time.Now(), false)

	id1 := "fb2211c9-5d53-4a44-895b-79c42174d521"
	todoDTO1 := dto.NewFactory().NewTodo(id1, userID, "things todo", false, time.Now(), time.Now(), false)

	s.Repository.On("FetchAllByUser", ctx, userDTO).Return([]*dto.Todo{todoDTO0, todoDTO1}, nil)

	// assert
	res, err := s.Usecase.FetchAllByUser(ctx, userDTO)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res, 2)
	assert.Equal(s.T(), todoDTO0.UserID, res[0].UserID)
	assert.Equal(s.T(), todoDTO0.UserID, res[1].UserID)
}

func (s *TodoServiceTestSuite) TestUpdateSuccess() {
	ctx := context.Background()

	// mock repo
	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	userDTO := dto.NewFactory().NewUser(userID, "account@emai.com", "strong-password", time.Now())

	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	todoDTO := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	s.Repository.On("Update", ctx, todoDTO).Return(nil)

	// assert
	res, err := s.Usecase.Update(ctx, userDTO, todoDTO)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), todoDTO.UserID, res.UserID)
}

func (s *TodoServiceTestSuite) TestDeleteSuccess() {
	ctx := context.Background()

	// mock repo
	userID := "2192fc7b-bd9b-446d-a50e-5ce0ba02cee6"
	userDTO := dto.NewFactory().NewUser(userID, "account@emai.com", "strong-password", time.Now())

	id := "4daaaea8-4721-4644-aaac-7958805b4530"
	todoDTO := dto.NewFactory().NewTodo(id, userID, "things todo", false, time.Now(), time.Now(), false)

	s.Repository.On("Update", ctx, todoDTO).Return(nil)

	// assert
	err := s.Usecase.Delete(ctx, userDTO, todoDTO)
	assert.NoError(s.T(), err)
}

func TestTodoService(t *testing.T) {
	suite.Run(t, new(TodoServiceTestSuite))
}

package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	Usecase    Usecase
	Repository *mocks.Repository
}

func (s *UserServiceTestSuite) SetupTest() {
	s.Repository = new(mocks.Repository)

	usecase, err := NewService(WithRepository(s.Repository))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create usecase: %s", err))
	}

	s.Usecase = usecase
}

func (s *UserServiceTestSuite) TestFailWhenEmailAlreadyExist() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("existing@mail.com", "PASSWORD")

	// mock repo
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(nil, nil)

	// assert
	_, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidSignUpReq)
}

func (s *UserServiceTestSuite) TestFailWhenTooShortPassword() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("existing@mail.com", "123")

	// assert
	_, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidSignUpReq)
	assert.Regexp(s.T(), "Error:Field validation", err)
}

func (s *UserServiceTestSuite) TestFailWhenInvalidRequest() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("invalid-email", "PASSWORD")

	// assert
	_, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidSignUpReq)
	assert.Regexp(s.T(), "Error:Field validation", err)
}

func (s *UserServiceTestSuite) TestSuccessWhenValidRequest() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("good-guy@mail.com", "STRONG-PASSWORD")

	// mock repo
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(nil, ErrNotFound)
	s.Repository.On("Store", ctx, mock.AnythingOfType("*dto.User")).Return(nil)

	// assert
	resp, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), req.Email, resp.Email)

}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

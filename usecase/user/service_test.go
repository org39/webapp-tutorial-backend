package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/pkg/crypt"
	auth_mocks "github.com/org39/webapp-tutorial-backend/usecase/auth/mocks"
	"github.com/org39/webapp-tutorial-backend/usecase/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	Usecase     Usecase
	AuthUsecase *auth_mocks.Usecase
	Repository  *mocks.Repository
}

func (s *UserServiceTestSuite) SetupTest() {
	s.Repository = new(mocks.Repository)
	s.AuthUsecase = new(auth_mocks.Usecase)

	usecase, err := NewService(
		WithRepository(s.Repository),
		WithAuthUsecase(s.AuthUsecase),
	)
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create usecase: %s", err))
	}

	s.Usecase = usecase
}

func (s *UserServiceTestSuite) TestSignUpFailWhenEmailAlreadyExist() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("existing@mail.com", "PASSWORD")

	// mock repo
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(nil, nil)

	// assert
	_, _, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenDatabaseError() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("valid@mail.com", "PASSWORD")

	// mock repo
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(nil, ErrDatabaseError)

	// assert
	_, _, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrDatabaseError)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenTooShortPassword() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("existing@mail.com", "123")

	// assert
	_, _, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
	assert.Regexp(s.T(), "Error:Field validation", err)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenInvalidRequest() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("invalid-email", "PASSWORD")

	// assert
	_, _, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
	assert.Regexp(s.T(), "Error:Field validation", err)
}

func (s *UserServiceTestSuite) TestSignUpSuccessWhenValidRequest() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserSignUpRequest("good-guy@mail.com", "STRONG-PASSWORD")

	// mock repo
	dummyToken := dto.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(nil, ErrNotFound)
	s.Repository.On("Store", ctx, mock.AnythingOfType("*dto.User")).Return(nil)
	s.AuthUsecase.On("GenereateToken", ctx, mock.AnythingOfType("*dto.AuthGenerateRequest")).Return(dummyToken, nil)

	// assert
	resp, tokens, err := s.Usecase.SignUp(ctx, req)
	s.Repository.AssertExpectations(s.T())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), req.Email, resp.Email)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func (s *UserServiceTestSuite) TestLoginSuccessWhenCorrectPassword() {
	ctx := context.Background()

	email := "good-guy@mail.com"
	plainPassword := "STRONG-PASSWORD"
	password, err := crypt.Hash([]byte(plainPassword))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to hash plainPassword: %s", err))
	}

	req := dto.NewFactory().NewUserLoginRequest(email, plainPassword)
	userDTO := dto.NewFactory().NewUser("id", email, password, time.Now())

	dummyToken := dto.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(userDTO, nil)
	s.AuthUsecase.On("GenereateToken", ctx, mock.AnythingOfType("*dto.AuthGenerateRequest")).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Login(ctx, req)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func (s *UserServiceTestSuite) TestLoginFailWhenWrongPassword() {
	ctx := context.Background()

	email := "good-guy@mail.com"
	plainPassword := "STRONG-PASSWORD"
	password, err := crypt.Hash([]byte(plainPassword))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to hash plainPassword: %s", err))
	}

	req := dto.NewFactory().NewUserLoginRequest(email, "WRONG-PASSWORD")
	userDTO := dto.NewFactory().NewUser("id", email, password, time.Now())

	dummyToken := dto.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, req.Email).Return(userDTO, nil)
	s.AuthUsecase.On("GenereateToken", ctx, mock.AnythingOfType("*dto.AuthGenerateRequest")).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Login(ctx, req)
	assert.ErrorIs(s.T(), err, ErrUnauthorized)
	assert.Nil(s.T(), tokens)
}

func (s *UserServiceTestSuite) TestRefreshSuccessWithValidToken() {
	ctx := context.Background()
	req := dto.NewFactory().NewUserRefreshRequest("VALID-TOKEN")

	dummyToken := dto.NewFactory().NewAuthTokenPair("access", "refresh")
	s.AuthUsecase.On("RefreshToken", ctx, mock.AnythingOfType("*dto.AuthRefreshRequest")).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Refresh(ctx, req)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

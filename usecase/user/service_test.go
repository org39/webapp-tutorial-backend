package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity"
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
	email := "existing@mail.com"
	password := "PASSWORD"

	// mock repo
	s.Repository.On("FetchByEmail", ctx, email).Return(nil, nil)

	// assert
	_, _, err := s.Usecase.SignUp(ctx, email, password)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenDatabaseError() {
	ctx := context.Background()
	email := "valid@mail.com"
	password := "PASSWORD"

	// mock repo
	s.Repository.On("FetchByEmail", ctx, email).Return(nil, ErrDatabaseError)

	// assert
	_, _, err := s.Usecase.SignUp(ctx, email, password)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrDatabaseError)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenTooShortPassword() {
	ctx := context.Background()
	email := "existing@mail.com"
	password := "PASS"

	// assert
	_, _, err := s.Usecase.SignUp(ctx, email, password)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
}

func (s *UserServiceTestSuite) TestSignUpFailWhenInvalidRequest() {
	ctx := context.Background()
	email := "invalid-email"
	password := "PASSWORD"

	// assert
	_, _, err := s.Usecase.SignUp(ctx, email, password)
	s.Repository.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
}

func (s *UserServiceTestSuite) TestSignUpSuccessWhenValidRequest() {
	ctx := context.Background()
	email := "good-guy@mail.com"
	password := "STRONG-PASSWORD"

	// mock repo
	dummyToken := entity.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, email).Return(nil, ErrNotFound)
	s.Repository.On("Store", ctx, mock.AnythingOfType("*dto.User")).Return(nil)
	s.AuthUsecase.On("GenereateToken", ctx, mock.AnythingOfType("string")).Return(dummyToken, nil)

	// assert
	resp, tokens, err := s.Usecase.SignUp(ctx, email, password)
	s.Repository.AssertExpectations(s.T())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), email, resp.Email)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func (s *UserServiceTestSuite) TestLoginSuccessWhenCorrectPassword() {
	ctx := context.Background()

	uuid := "62db52ec-5c8a-4a3c-a3c4-0b69db9a1f30"
	email := "good-guy@mail.com"
	plainPassword := "STRONG-PASSWORD"
	password, err := crypt.Hash([]byte(plainPassword))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to hash plainPassword: %s", err))
	}

	userDTO := dto.NewFactory().NewUser(uuid, email, password, time.Now())

	dummyToken := entity.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, email).Return(userDTO, nil)
	s.AuthUsecase.On("GenereateToken", ctx, uuid).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Login(ctx, email, plainPassword)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func (s *UserServiceTestSuite) TestLoginFailWhenWrongPassword() {
	ctx := context.Background()

	uuid := "62db52ec-5c8a-4a3c-a3c4-0b69db9a1f30"
	email := "good-guy@mail.com"
	plainPassword := "STRONG-PASSWORD"
	password, err := crypt.Hash([]byte(plainPassword))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to hash plainPassword: %s", err))
	}

	userDTO := dto.NewFactory().NewUser(uuid, email, password, time.Now())

	dummyToken := entity.NewFactory().NewAuthTokenPair("access", "refresh")
	s.Repository.On("FetchByEmail", ctx, email).Return(userDTO, nil)
	s.AuthUsecase.On("GenereateToken", ctx, uuid).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Login(ctx, email, "WRONG-PASSWORD")
	assert.ErrorIs(s.T(), err, ErrUnauthorized)
	assert.Nil(s.T(), tokens)
}

func (s *UserServiceTestSuite) TestRefreshSuccessWithValidToken() {
	ctx := context.Background()
	refreshToken := "VALID-TOKEN"

	dummyToken := entity.NewFactory().NewAuthTokenPair("access", "refresh")
	s.AuthUsecase.On("RefreshToken", ctx, refreshToken).Return(dummyToken, nil)

	// assert
	tokens, err := s.Usecase.Refresh(ctx, refreshToken)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokens.AccessToken)
	assert.NotEmpty(s.T(), tokens.RefreshToken)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

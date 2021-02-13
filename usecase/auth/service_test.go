package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/org39/webapp-tutorial-backend/entity/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	Usecase Usecase
}

func (s *AuthServiceTestSuite) SetupTest() {
	usecase, err := NewService(WithSecret("top-secret"))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create usecase: %s", err))
	}

	s.Usecase = usecase
}

func (s *AuthServiceTestSuite) TestSuccessGenereateToken() {
	ctx := context.Background()
	req := dto.NewFactory().NewAuthGenerateRequest("hatsune@miku.com")

	// assert
	tokenPair, err := s.Usecase.GenereateToken(ctx, req)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokenPair.AccessToken)
	assert.NotEmpty(s.T(), tokenPair.RefreshToken)
}

func (s *AuthServiceTestSuite) TestFailGenereateTokenWithInvalidEmail() {
	ctx := context.Background()
	req := dto.NewFactory().NewAuthGenerateRequest("invalid-email")

	// assert
	tokenPair, err := s.Usecase.GenereateToken(ctx, req)
	assert.ErrorIs(s.T(), err, ErrInvalidRequest)
	assert.Nil(s.T(), tokenPair)
}

func (s *AuthServiceTestSuite) TestSuccessRefreshWithValidToken() {
	ctx := context.Background()
	tokenReq := dto.NewFactory().NewAuthGenerateRequest("hatsune@miku.com")

	// assert
	tokenPair, err := s.Usecase.GenereateToken(ctx, tokenReq)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokenPair.AccessToken)
	assert.NotEmpty(s.T(), tokenPair.RefreshToken)

	refreshReq := dto.NewFactory().NewAuthRefreshRequest(tokenPair.RefreshToken)
	newTokenPair, err := s.Usecase.RefreshToken(ctx, refreshReq)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), newTokenPair.AccessToken)
	assert.NotEmpty(s.T(), newTokenPair.RefreshToken)
}

func (s *AuthServiceTestSuite) TestSuccessVerifyWithValidToken() {
	ctx := context.Background()
	tokenReq := dto.NewFactory().NewAuthGenerateRequest("hatsune@miku.com")

	// assert
	tokenPair, err := s.Usecase.GenereateToken(ctx, tokenReq)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), tokenPair.AccessToken)
	assert.NotEmpty(s.T(), tokenPair.RefreshToken)

	verifyReq := dto.NewFactory().NewAuthVerifyRequest(tokenPair.AccessToken)
	err = s.Usecase.VerifyToken(ctx, verifyReq)
	assert.NoError(s.T(), err)
}

func (s *AuthServiceTestSuite) TestFailVerifyWithInvalidToken() {
	ctx := context.Background()
	verifyReq := dto.NewFactory().NewAuthVerifyRequest("invalid-token")
	err := s.Usecase.VerifyToken(ctx, verifyReq)
	assert.ErrorIs(s.T(), err, ErrUnauthorized)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

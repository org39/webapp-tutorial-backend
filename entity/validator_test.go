package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityValidatorSuite struct {
	suite.Suite
	Validator *Validator
}

func (s *EntityValidatorSuite) SetupTest() {
	s.Validator = NewValidator()
}

func (s *EntityValidatorSuite) TestEmailValidatorSuccess() {
	cases := []struct {
		email string
	}{
		{email: "test@example.com"},
	}

	for _, c := range cases {
		err := s.Validator.ValidateEmail(c.email)
		assert.NoError(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestEmailValidatorFailure() {
	cases := []struct {
		email string
	}{
		{email: "tes_example.com"},
	}

	for _, c := range cases {
		err := s.Validator.ValidateEmail(c.email)
		assert.Error(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestPlainPasswordSuccess() {
	cases := []struct {
		password string
	}{
		{password: "STRONG_PASSWORD"},
	}

	for _, c := range cases {
		err := s.Validator.ValidatePlainPassword(c.password)
		assert.NoError(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestPlainPasswordFailure() {
	cases := []struct {
		password string
	}{
		{password: "AAA"},
	}

	for _, c := range cases {
		err := s.Validator.ValidatePlainPassword(c.password)
		assert.Error(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestIDSuccess() {
	cases := []struct {
		id string
	}{
		{id: "6ce768ac-a944-4268-a51f-bed2bb551cb5"},
	}

	for _, c := range cases {
		err := s.Validator.ValidateID(c.id)
		assert.NoError(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestIDFailure() {
	cases := []struct {
		id string
	}{
		{id: "wrong-id"},
	}

	for _, c := range cases {
		err := s.Validator.ValidateID(c.id)
		assert.Error(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestTokenSuccess() {
	cases := []struct {
		token string
	}{
		{token: "VALID-TOKEN"},
	}

	for _, c := range cases {
		err := s.Validator.ValidateToken(c.token)
		assert.NoError(s.T(), err)
	}
}

func (s *EntityValidatorSuite) TestTokenFailure() {
	cases := []struct {
		token string
	}{
		{token: ""},
	}

	for _, c := range cases {
		err := s.Validator.ValidateToken(c.token)
		assert.Error(s.T(), err)
	}
}

func TestEntityValidator(t *testing.T) {
	suite.Run(t, new(EntityValidatorSuite))
}

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityUserTestSuite struct {
	suite.Suite
}

func (s *EntityUserTestSuite) TestUserInvalid() {
	cases := []struct {
		email    string
		password string
	}{
		{email: "", password: ""},
		{email: "invalid", password: "PASSWORD"},
	}

	for _, c := range cases {
		badass, err := NewFactory().NewUser(c.email, c.password)
		assert.NoError(s.T(), err)

		v := badass.Valid()
		assert.Error(s.T(), v)
	}
}

func TestEntityUser(t *testing.T) {
	suite.Run(t, new(EntityUserTestSuite))
}

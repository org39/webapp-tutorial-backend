package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityUserTestSuite struct {
	suite.Suite
}

func (s *EntityUserTestSuite) SetupTest() {
}

func (s *EntityUserTestSuite) TestUserInvalid() {
	cases := []struct {
		email     string
		password  string
		createdAt time.Time
	}{
		{email: "", password: "", createdAt: time.Now()},
		{email: "invalid", password: "PASSWORD", createdAt: time.Now()},
	}

	for _, c := range cases {
		badass, err := NewFactory().NewUser(c.email, c.password, c.createdAt)
		assert.NoError(s.T(), err)

		v := badass.Valid()
		assert.Error(s.T(), v)
	}
}

func TestEntityUser(t *testing.T) {
	suite.Run(t, new(EntityUserTestSuite))
}

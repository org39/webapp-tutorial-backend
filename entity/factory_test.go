package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityFactoryTestSuite struct {
	suite.Suite
	Factory *Factory
}

func (s *EntityFactoryTestSuite) SetupTest() {
	s.Factory = NewFactory()
}

func (s *EntityFactoryTestSuite) TestCreateValidUser() {
	cases := []struct {
		email    string
		password string
	}{
		{email: "hatsune@miku.com", password: "PASSWORD"},
	}

	for _, c := range cases {
		good, err := s.Factory.NewUser(c.email, c.password)
		assert.NoError(s.T(), err)

		v := good.Valid()
		assert.NoError(s.T(), v)
	}
}

func (s *EntityFactoryTestSuite) TestCreateValidTodo() {
	u, err := NewFactory().NewUser("hatsune@miku", "PASSWORD")
	assert.NoError(s.T(), err)

	cases := []struct {
		user    *User
		content string
	}{
		{user: u, content: "THINGS TODO"},
	}

	for _, c := range cases {
		good, err := s.Factory.NewTodo(c.user, c.content)
		assert.NoError(s.T(), err)

		v := good.Valid()
		assert.NoError(s.T(), v)
	}
}

func TestEntityFactory(t *testing.T) {
	suite.Run(t, new(EntityFactoryTestSuite))
}

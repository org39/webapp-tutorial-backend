package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityTodoTestSuite struct {
	suite.Suite
}

func (s *EntityTodoTestSuite) TestCreationValid() {
	u, err := NewFactory().NewUser("hatsnune@miku.com", "very-strong-password")
	assert.NoError(s.T(), err)

	cases := []struct {
		user    *User
		content string
	}{
		{user: u, content: "TODO1"},
	}

	for _, c := range cases {
		e, err := NewFactory().NewTodo(u, c.content)
		assert.NoError(s.T(), err)

		v := e.Valid()
		assert.NoError(s.T(), v)
	}
}

func TestEntityTodo(t *testing.T) {
	suite.Run(t, new(EntityTodoTestSuite))
}

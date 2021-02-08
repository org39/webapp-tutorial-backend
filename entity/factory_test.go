package entity

import (
	"testing"
	"time"

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
		email     string
		password  string
		createdAt time.Time
	}{
		{email: "hatsune@miku.com", password: "PASSWORD", createdAt: time.Now()},
	}

	for _, c := range cases {
		good, err := s.Factory.NewUser(c.email, c.password, c.createdAt)
		assert.NoError(s.T(), err)

		v := good.Valid()
		assert.NoError(s.T(), v)
	}
}

func TestEntityFactory(t *testing.T) {
	suite.Run(t, new(EntityFactoryTestSuite))
}

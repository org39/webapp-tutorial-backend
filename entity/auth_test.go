package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityAuthTestSuite struct {
	suite.Suite
}

func (t *EntityAuthTestSuite) TestCreationValid() {
	tokens := NewFactory().NewAuthTokenPair("token", "refresh_token")
	err := tokens.Valid()
	assert.NoError(t.T(), err)
}

func TestEntityAuth(t *testing.T) {
	suite.Run(t, new(EntityAuthTestSuite))
}

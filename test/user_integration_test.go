package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/org39/webapp-tutorial-backend/app"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	jpassert "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	Application *app.App
	Server      *echo.Echo
}

func (s *AuthIntegrationTestSuite) SetupTest() {
	application, server, err := buildTestServer()
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create test Server: %s", err))
	}

	s.Application = application
	s.Server = server

	_, err = s.Application.DB.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", s.Application.Config.UserTable))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to truncate %s table: %s", s.Application.Config.UserTable, err))
	}
}

func (s *AuthIntegrationTestSuite) TearDownTest() {
	s.Application.DB.Close()
	app.ClearDepencencyGraph()
}

func (s *AuthIntegrationTestSuite) apiTest(name string) *apitest.APITest {
	return apitest.New(name).
		Recorder(recorder).
		Report(apitest.SequenceDiagram(fmt.Sprintf("apitest/%s", name))).
		Handler(s.Server)
}

func (s *AuthIntegrationTestSuite) TestRegisterSuccess() {
	email := "hatsune@miku.com"
	password := "very-strong-password"

	s.apiTest("user/register").
		Post("/user/register").
		JSON(map[string]string{
			"email":    email,
			"password": password,
		}).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Status(http.StatusCreated).
		End()
}

func (s *AuthIntegrationTestSuite) TestLoginRefreshSuccess() {
	email := "hatsune@miku.com"
	password := "very-strong-password"

	s.apiTest("user/login").
		Post("/user/register").
		JSON(map[string]string{
			"email":    email,
			"password": password,
		}).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Status(http.StatusCreated).
		End()

	s.apiTest("user/login").
		Post("/user/login").
		JSON(map[string]string{
			"email":    email,
			"password": password,
		}).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Assert(jpassert.Present("$.access_token")).
		Status(http.StatusOK).
		End()
}

func TestAuthIntegrationTest(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}

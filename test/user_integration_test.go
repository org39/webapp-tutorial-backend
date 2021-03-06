package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	app "github.com/org39/webapp-tutorial-backend/app/server"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/testreport"
	"github.com/steinfletcher/apitest"
	jpassert "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserIntegrationTestSuite struct {
	suite.Suite

	Application       *app.App
	Server            *echo.Echo
	TestSuiteReporter *testreport.TestSuiteReporter
}

func (s *UserIntegrationTestSuite) SetupSuite() {
	reporter := testreport.New("UserIntegrationTest", "./report")
	application, server, err := buildTestServer()
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create test Server: %s", err))
	}

	s.Application = application
	s.Server = server
	s.TestSuiteReporter = reporter
}

func (s *UserIntegrationTestSuite) SetupTest() {
	_, err := s.Application.DB.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", s.Application.Config.UserTable))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to truncate %s table: %s", s.Application.Config.UserTable, err))
	}
}

func (s *UserIntegrationTestSuite) TearDownSuite() {
	s.Application.DB.Close()
	s.TestSuiteReporter.Flush()
}

func (s *UserIntegrationTestSuite) apiTest(name string) *apitest.APITest {
	return apitest.New(name).
		Recorder(recorder).
		Report(s.TestSuiteReporter).
		Handler(s.Server)
}

func (s *UserIntegrationTestSuite) TestRegisterSuccess() {
	createTestAccount(s.T(), s.apiTest("TestRegisterSuccess"))
}

func (s *UserIntegrationTestSuite) TestLoginRefreshSuccess() {
	account := createTestAccount(s.T(), s.apiTest("TestLoginRefreshSuccess"))
	loginResp := s.apiTest("TestLoginRefreshSuccess").
		Post("/user/login").
		JSON(map[string]string{
			"email":    account.User.Email,
			"password": account.Password,
		}).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Assert(jpassert.Present("$.access_token")).
		Status(http.StatusOK).
		End().Response

	refreshTokenCookie := findCookieByName(loginResp.Cookies(), "refresh_token")
	s.apiTest("TestLoginRefreshSuccess").
		Post("/user/refresh").
		Cookie("refresh_token", refreshTokenCookie.Value).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Assert(jpassert.Present("$.access_token")).
		Status(http.StatusOK).
		End()
}

func (s *UserIntegrationTestSuite) TestGetUserSuccess() {
	account := createTestAccount(s.T(), s.apiTest("TestGetUserSuccess"))
	s.apiTest("TestGetUserSuccess").
		Post("/user/login").
		JSON(map[string]string{
			"email":    account.User.Email,
			"password": account.Password,
		}).
		Expect(s.T()).
		CookiePresent("refresh_token").
		Assert(jpassert.Present("$.access_token")).
		Status(http.StatusOK).
		End()
	s.apiTest("TestGetUserSuccess").
		Get("/user").
		Header("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)).
		Expect(s.T()).
		Assert(jpassert.Equal("$.email", account.User.Email)).
		Status(http.StatusOK).
		End()
}

func TestUserIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	suite.Run(t, new(UserIntegrationTestSuite))
}

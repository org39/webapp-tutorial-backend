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

type TodoIntegrationTestSuite struct {
	suite.Suite

	Application       *app.App
	Server            *echo.Echo
	TestSuiteReporter *testreport.TestSuiteReporter
}

func (s *TodoIntegrationTestSuite) SetupSuite() {
	reporter := testreport.New("TodoIntegrationTest", "./report")
	application, server, err := buildTestServer()
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to create test Server: %s", err))
	}

	s.Application = application
	s.Server = server
	s.TestSuiteReporter = reporter
}

func (s *TodoIntegrationTestSuite) SetupTest() {
	_, err := s.Application.DB.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", s.Application.Config.UserTable))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to truncate %s table: %s", s.Application.Config.UserTable, err))
	}

	_, err = s.Application.DB.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", s.Application.Config.TodoTable))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("fail to truncate %s table: %s", s.Application.Config.TodoTable, err))
	}
}

func (s *TodoIntegrationTestSuite) TearDownSuite() {
	s.Application.DB.Close()
	s.TestSuiteReporter.Flush()
}

func (s *TodoIntegrationTestSuite) apiTest(name string) *apitest.APITest {
	return apitest.New(name).
		Recorder(recorder).
		Report(s.TestSuiteReporter).
		Handler(s.Server)
}

func (s *TodoIntegrationTestSuite) TestCreateTodoSuccess() {
	content := "things todo"
	account := createTestAccount(s.T(), s.apiTest("TestCreateTodoSuccess"))
	todo := createTestTodo(s.T(), s.apiTest("TestCreateTodoSuccess"), account, content)

	assert.Equal(s.T(), content, todo.Content)
	assert.False(s.T(), todo.Completed)
	assert.False(s.T(), todo.Deleted)
}

func (s *TodoIntegrationTestSuite) TestMarkCompletedSuccess() {
	content := "things todo"
	account := createTestAccount(s.T(), s.apiTest("TestMarkCompletedSuccess"))
	todo := createTestTodo(s.T(), s.apiTest("TestMarkCompletedSuccess"), account, content)

	resourceURL := fmt.Sprintf("/todos/%s", todo.ID)
	newContent := "Good Good"
	completed := true
	s.apiTest("TestMarkCompletedSuccess").
		Put(resourceURL).
		Header("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)).
		JSON(map[string]interface{}{
			"id":        todo.ID,
			"user_id":   account.User.ID,
			"content":   newContent,
			"completed": completed,
			"deleted":   false,
		}).
		Expect(s.T()).
		Assert(jpassert.Equal("$.content", newContent)).
		Assert(jpassert.Equal("$.completed", completed)).
		Assert(jpassert.Equal("$.deleted", false)).
		Status(http.StatusOK).
		End()
}

func (s *TodoIntegrationTestSuite) TestGetAllTodosSuccess() {
	content := "things todo"
	account := createTestAccount(s.T(), s.apiTest("TestGetAllTodosSuccess"))
	_ = createTestTodo(s.T(), s.apiTest("TestGetAllTodosSuccess"), account, content)

	s.apiTest("TestGetAllTodosSuccess").
		Get("/todos").
		Header("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)).
		Expect(s.T()).
		Assert(jpassert.Len("$", 1)).
		Assert(jpassert.Equal("$[0].content", content)).
		Assert(jpassert.Equal("$[0].completed", false)).
		Assert(jpassert.Equal("$[0].deleted", false)).
		Status(http.StatusOK).
		End()
}

func (s *TodoIntegrationTestSuite) TestGetTodoByIdSuccess() {
	content := "things todo"
	account := createTestAccount(s.T(), s.apiTest("TestGetTodoByIdSuccess"))
	todo := createTestTodo(s.T(), s.apiTest("TestGetTodoByIdSuccess"), account, content)

	resourceURL := fmt.Sprintf("/todos/%s", todo.ID)
	s.apiTest("TestGetTodoByIdSuccess").
		Get(resourceURL).
		Header("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)).
		Expect(s.T()).
		Assert(jpassert.Equal("$.id", todo.ID)).
		Assert(jpassert.Equal("$.content", todo.Content)).
		Assert(jpassert.Equal("$.completed", todo.Completed)).
		Assert(jpassert.Equal("$.deleted", todo.Deleted)).
		Status(http.StatusOK).
		End()
}

func TestTodoIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	suite.Run(t, new(TodoIntegrationTestSuite))
}

package test

import (
	"database/sql/driver"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	jpassert "github.com/steinfletcher/apitest-jsonpath"
	apitestdb "github.com/steinfletcher/apitest/x/db"
	"github.com/stretchr/testify/assert"
)

var recorder *apitest.Recorder

func init() {
	recorder = apitest.NewTestRecorder()
}

func buildTestServer() (*app.App, *echo.Echo, error) {
	// clear dependency graph
	app.ClearDepencencyGraph()

	// build application
	application, err := app.New(newRecorededMysqlConn)
	if err != nil {
		return nil, nil, err
	}

	// disable log
	application.RootLogger.Logger.SetLevel(logrus.PanicLevel)

	// attach application to RestAPI presenter
	server, err := rest.New(&app.DepencencyInjector, application.RootLogger, func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err != nil {
		return nil, nil, err
	}

	return application, server, nil
}

func newRecorededMysqlConn(conf *app.Config) (driver.Connector, error) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	databaseHost := net.JoinHostPort(conf.DatabaseHost, conf.DatabasePort)
	dsn := &mysql.Config{
		Addr:                 databaseHost,
		Net:                  "tcp",
		User:                 conf.DatabaseUser,
		Passwd:               conf.DatabasePass,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  utc,
		ParseTime:            true,
		DBName:               conf.DatabaseName,
		AllowNativePasswords: true,
	}

	con, err := mysql.NewConnector(dsn)
	if err != nil {
		return nil, err
	}

	return apitestdb.WrapConnectorWithRecorder(con, "mysql", recorder), nil
}

func findCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
type Account struct {
	User         User `json:"user"`
	Password     string
	AccessToken  string `json:"access_token"`
	RefreshToken string
}

func createTestAccount(t *testing.T, apiTest *apitest.APITest) Account {
	account := Account{
		User: User{
			Email: "hatsune@miku.com",
		},
		Password: "very-strong-password",
	}

	res := apiTest.Post("/user/register").
		JSON(map[string]string{
			"email":    account.User.Email,
			"password": account.Password,
		}).
		Expect(t).
		CookiePresent("refresh_token").
		Assert(jpassert.Present("$.access_token")).
		Status(http.StatusCreated).
		End()

	// fetch accessToken from response body
	res.JSON(&account)

	// fetch refreshToken from cookies
	resp := res.Response
	loginRespCookies := resp.Cookies()
	refreshTokenCookie := findCookieByName(loginRespCookies, "refresh_token")
	assert.NotNil(t, refreshTokenCookie)
	account.RefreshToken = refreshTokenCookie.Value

	return account
}

type Todo struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Deleted   bool   `json:"deleted"`
}

func createTestTodo(t *testing.T, apiTest *apitest.APITest, account Account, content string) Todo {
	res := apiTest.Post("/todos").
		JSON(map[string]string{
			"content": content,
		}).
		Header("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)).
		Expect(t).
		Assert(jpassert.Equal("$.content", content)).
		Assert(jpassert.Equal("$.completed", false)).
		Assert(jpassert.Equal("$.deleted", false)).
		Status(http.StatusCreated).
		End()

	// fetch newly created todo from response body
	todo := Todo{}
	res.JSON(&todo)

	return todo
}

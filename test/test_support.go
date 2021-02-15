package test

import (
	"database/sql/driver"
	"net"
	"net/http"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	apitestdb "github.com/steinfletcher/apitest/x/db"
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

package test

import (
	"database/sql/driver"
	"net"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/facebookgo/inject"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/router"
	"github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	apitestdb "github.com/steinfletcher/apitest/x/db"
)

var recorder *apitest.Recorder

func init() {
	recorder = apitest.NewTestRecorder()
}

func buildTestServer() (*app.App, *echo.Echo, error) {
	// build application
	application, err := app.New(newRecorededMysqlConn)
	if err != nil {
		return nil, nil, err
	}

	// disable log
	application.RootLogger.Logger.SetLevel(logrus.PanicLevel)

	// attach application to RestAPI presenter
	server, err := newRestAPI(application)
	if err != nil {
		return nil, nil, err
	}

	return application, server, nil
}

func newRestAPI(application *app.App) (*echo.Echo, error) {
	server, err := router.New(application.RootLogger)
	if err != nil {
		return nil, err
	}

	restAPI, err := rest.NewDispatcher()
	if err != nil {
		return nil, err
	}

	// middleware
	authm := new(rest.AuthMiddleware)
	if err := app.DepencencyInjector.Provide(&inject.Object{Value: authm}); err != nil {
		return nil, err
	}

	// user RestAPI
	userAPI := new(rest.UserDispatcher)
	restAPI.AttachDispatcher(userAPI)
	if err := app.DepencencyInjector.Provide(&inject.Object{Value: userAPI}); err != nil {
		return nil, err
	}

	// build dependency graph
	if err := app.DepencencyInjector.Populate(); err != nil {
		return nil, err
	}

	restAPI.Dispatch(server)
	return server, nil
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

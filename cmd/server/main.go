package main

import (
	"context"
	"database/sql/driver"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	// build application
	application, err := app.New(newMysqlConn)
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()

	// attach application to RestAPI presenter
	server, err := rest.New(&app.DepencencyInjector, application.RootLogger, readiness(application))
	if err != nil {
		panic(err)
	}

	// server start and wait signal or error
	quit := make(chan os.Signal, 5)
	signal.Notify(quit, os.Interrupt)
	serverErr := make(chan error)
	go func() {
		serverErr <- server.Start(":8080")
	}()

	// wait signal or error
	select {
	case <-quit:
		application.RootLogger.Info("receive signal")
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			application.RootLogger.WithField("error", err).Fatal("server shutdown error")
		} else {
			application.RootLogger.Info("stop server")
		}
	case err := <-serverErr:
		application.RootLogger.WithField("error", err).Fatal("server start error")
	}
}

func readiness(application *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := application.DB.Ping(); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func newMysqlConn(conf *app.Config) (driver.Connector, error) {
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

	return mysql.NewConnector(dsn)
}

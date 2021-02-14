package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/facebookgo/inject"
	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/router"
)

func main() {
	// build application
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()

	// attach application to RestAPI presenter
	server, err := newRestAPI(application)
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

func newRestAPI(application *app.App) (*echo.Echo, error) {
	server, err := router.New(application.RootLogger)
	if err != nil {
		return nil, err
	}

	restAPI, err := rest.NewDispatcher(
		rest.WithReadinessCheck(readiness(application)),
	)
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

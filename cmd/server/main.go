package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/org39/webapp-tutorial-backend/pkg/router"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}

	server, err := router.New(application.RootLogger)
	if err != nil {
		panic(err)
	}

	restAPI, err := rest.NewDispatcher(
		rest.WithUserDispatcher(application.UserUsecase, application.RootLogger),
	)
	if err != nil {
		panic(err)
	}

	restAPI.Dispatch(server)

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

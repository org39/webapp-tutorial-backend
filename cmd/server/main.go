package main

import (
	"github.com/org39/webapp-tutorial-backend/app"
	"github.com/org39/webapp-tutorial-backend/presenter/rest"

	"github.com/org39/webapp-tutorial-backend/pkg/router"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}

	restRouter, err := router.New(application.RootLogger)
	if err != nil {
		panic(err)
	}

	restAPI, err := rest.NewDispatcher(
		rest.WithUserDispatcher(application.UserUsecase, application.RootLogger),
	)
	if err != nil {
		panic(err)
	}

	application.RootLogger.WithField("app", application.UserUsecase).Debug("test")

	restAPI.Dispatch(restRouter)
	restRouter.Start(":8080")
}

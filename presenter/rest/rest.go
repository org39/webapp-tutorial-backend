package rest

import (
	"github.com/facebookgo/inject"
	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
	"github.com/org39/webapp-tutorial-backend/pkg/router"
)

func New(g *inject.Graph, logger *log.Logger, readiness echo.HandlerFunc) (*echo.Echo, error) {
	server, err := router.New(logger)
	if err != nil {
		return nil, err
	}

	restAPI, err := NewDispatcher(
		WithReadinessCheck(readiness),
	)
	if err != nil {
		return nil, err
	}

	// middleware
	authm := new(AuthMiddleware)
	if err := g.Provide(&inject.Object{Value: authm}); err != nil {
		return nil, err
	}

	// user RestAPI
	userAPI := new(UserDispatcher)
	restAPI.AttachDispatcher(userAPI)
	if err := g.Provide(&inject.Object{Value: userAPI}); err != nil {
		return nil, err
	}

	// todo RestAPI
	todoAPI := new(TodoDispatcher)
	restAPI.AttachDispatcher(todoAPI)
	if err := g.Provide(&inject.Object{Value: todoAPI}); err != nil {
		return nil, err
	}

	// build dependency graph
	if err := g.Populate(); err != nil {
		return nil, err
	}

	restAPI.Dispatch(server)
	return server, nil
}

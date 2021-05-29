package rest

import (
	"errors"
	"net/http"

	"github.com/org39/webapp-tutorial-backend/presenter/rest/rr"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

type TodoDispatcher struct {
	TodoUsecase    todo.Usecase    `inject:""`
	UserUsecase    user.Usecase    `inject:""`
	AuthMiddleware *AuthMiddleware `inject:""`
}

func (d *TodoDispatcher) Dispatch(e *echo.Echo) {
	auth := d.AuthMiddleware.Middleware()

	e.GET("todos", d.GetAllByUser(), auth)
	e.GET("todos/:id", d.GetByID(), auth)
	e.POST("todos", d.Create(), auth)
	e.PUT("todos/:id", d.UpdateByID(), auth)
	e.DELETE("todos/:id", d.DeleteByID(), auth)
}

func (d *TodoDispatcher) GetAllByUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := log.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		user, err := d.UserUsecase.FetchByID(ctx, authCtx.UserID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		todos, err := d.TodoUsecase.FetchAllByUser(ctx, user, false, false)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSON(http.StatusOK,
			rr.NewFactory().NewTodosResponse(todos),
		)
	}
}

func (d *TodoDispatcher) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := log.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		user, err := d.UserUsecase.FetchByID(ctx, authCtx.UserID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		payload, err := rr.NewFactory().NewTodoCreatRequest(c)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		todo, err := d.TodoUsecase.Create(ctx, user, payload.Content)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSON(http.StatusCreated,
			rr.NewFactory().NewTodoResponse(todo),
		)
	}
}

func (d *TodoDispatcher) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := log.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.UserID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		todo, err := d.TodoUsecase.FetchByID(ctx, userDTO, id)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSON(http.StatusOK,
			rr.NewFactory().NewTodoResponse(todo),
		)
	}
}

func (d *TodoDispatcher) UpdateByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := log.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		user, err := d.UserUsecase.FetchByID(ctx, authCtx.UserID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		payload, err := rr.NewFactory().NewTodoUpdateRequest(c)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		todo, err := d.TodoUsecase.Update(ctx, user, id, payload.Content, payload.Completed, payload.Deleted)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSON(http.StatusOK,
			rr.NewFactory().NewTodoResponse(todo),
		)
	}
}

func (d *TodoDispatcher) DeleteByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := log.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		user, err := d.UserUsecase.FetchByID(ctx, authCtx.UserID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		if err := d.TodoUsecase.Delete(ctx, user, id); err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func toTodoHTTPError(logger *log.Logger, err error) error {
	// errors defined in usecase
	switch {
	case errors.Is(err, todo.ErrInvalidRequest):
		return echo.NewHTTPError(http.StatusBadRequest)

	case errors.Is(err, todo.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound)

	case errors.Is(err, todo.ErrSystemError):
		logger.WithError(err).Error()
		return echo.NewHTTPError(http.StatusInternalServerError)

	case errors.Is(err, todo.ErrDatabaseError):
		logger.WithError(err).Error()
		return echo.NewHTTPError(http.StatusInternalServerError)

	case errors.Is(err, todo.ErrUnauthorized):
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	logger.WithError(err).Error()
	return echo.NewHTTPError(http.StatusInternalServerError)
}

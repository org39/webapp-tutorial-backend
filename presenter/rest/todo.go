package rest

import (
	"errors"
	"net/http"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

type TodoDispatcher struct {
	TodoUsecase todo.Usecase    `inject:""`
	UserUsecase user.Usecase    `inject:""`
	Authm       *AuthMiddleware `inject:""`
	Logger      *log.Logger     `inject:""`
}

func (d *TodoDispatcher) Dispatch(e *echo.Echo) {
	auth := d.Authm.Middleware()

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
		logger := d.Logger.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.ID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		todos, err := d.TodoUsecase.FetchAllByUser(ctx, userDTO)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSONPretty(http.StatusOK, todos, "  ")
	}
}

func (d *TodoDispatcher) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.ID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		payload := dto.NewFactory().NewTodoCreatRequest("")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		todoDTO, err := d.TodoUsecase.Create(ctx, userDTO, payload.Content)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSONPretty(http.StatusCreated, todoDTO, "  ")
	}
}

func (d *TodoDispatcher) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.ID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		todoDTO, err := d.TodoUsecase.FetchByID(ctx, userDTO, id)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSONPretty(http.StatusOK, todoDTO, "  ")
	}
}

func (d *TodoDispatcher) UpdateByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.ID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		payload := dto.NewFactory().NewTodoUpdateRequest("", false, false)
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		todoDTO, err := d.TodoUsecase.Update(ctx, userDTO, id, payload)
		if err != nil {
			return toTodoHTTPError(logger, err)
		}

		return c.JSONPretty(http.StatusOK, todoDTO, "  ")
	}
}

func (d *TodoDispatcher) DeleteByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		authCtx, ok := c.(*AuthorizedContext)
		if !ok {
			logger.WithError(errors.New("invalid authorized context")).Error()
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		userDTO, err := d.UserUsecase.FetchByID(ctx, authCtx.ID())
		if err != nil {
			return toHTTPError(logger, err)
		}

		id := c.Param("id")
		if err := d.TodoUsecase.Delete(ctx, userDTO, id); err != nil {
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

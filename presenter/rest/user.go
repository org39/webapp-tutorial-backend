package rest

import (
	// "context"
	"errors"
	// "fmt"
	// "time"
	"net/http"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

type userDispatcher struct {
	UserUsecase user.Usecase
	Logger      *log.Logger
}

func WithUserDispatcher(u user.Usecase, l *log.Logger) func(*Dispatcher) error {
	return func(r *Dispatcher) error {
		d := &userDispatcher{
			UserUsecase: u,
			Logger:      l,
		}
		r.dispatchers = append(r.dispatchers, d)
		return nil
	}
}

func (d *userDispatcher) Dispatch(e *echo.Echo) {
	e.POST("auth/register", d.Register())
}

func (d *userDispatcher) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		payload := dto.NewFactory().NewUserSignUpRequest("", "")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		response, err := d.UserUsecase.SignUp(ctx, payload)
		if err != nil {
			logger.WithField("error", err).Error("")
			return toHTTPError(err)
		}

		return c.JSONPretty(http.StatusCreated, response, "  ")
	}
}

func toHTTPError(err error) error {
	// errors defined in usecase
	switch {
	case errors.Is(err, user.ErrInvalidSignUpReq):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, user.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrSystemError):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	case errors.Is(err, user.ErrDatabaseError):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

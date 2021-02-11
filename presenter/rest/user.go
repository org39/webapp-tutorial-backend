package rest

import (
	"errors"
	"net/http"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/auth"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

type UserDispatcher struct {
	UserUsecase        user.Usecase `inject:""`
	AuthUsercase       auth.Usecase `inject:""`
	SecureRefreshToken bool         `inject:"rest.auth.secure_refresh_token"`
	Logger             *log.Logger  `inject:""`
}

func (d *UserDispatcher) Dispatch(e *echo.Echo) {
	e.POST("auth/register", d.Register())
}

func (d *UserDispatcher) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		payload := dto.NewFactory().NewUserSignUpRequest("", "")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		response, tokens, err := d.UserUsecase.SignUp(ctx, payload)
		if err != nil {
			logger.WithField("error", err).Error("")
			return toHTTPError(err)
		}

		// set refresh token as cookie
		cookie := new(http.Cookie)
		cookie.Name = "refresh_token"
		cookie.Value = tokens.RefreshToken
		if d.SecureRefreshToken {
			cookie.Secure = true
		}
		c.SetCookie(cookie)

		return c.JSONPretty(http.StatusCreated, map[string]interface{}{
			"user":         response,
			"access_token": tokens.AccessToken}, "  ")
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

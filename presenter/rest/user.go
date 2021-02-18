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

const (
	refreshTokenCookie = "refresh_token"
)

type UserDispatcher struct {
	UserUsecase        user.Usecase `inject:""`
	AuthUsercase       auth.Usecase `inject:""`
	SecureRefreshToken bool         `inject:"rest.auth.secure_refresh_token"`
	Logger             *log.Logger  `inject:""`
}

func (d *UserDispatcher) Dispatch(e *echo.Echo) {
	e.POST("user/register", d.Register())
	e.POST("user/login", d.Login())
	e.POST("user/refresh", d.Refresh())
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
			return toHTTPError(logger, err)
		}

		// set refresh token as cookie
		cookie := new(http.Cookie)
		cookie.Name = refreshTokenCookie
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

func (d *UserDispatcher) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		payload := dto.NewFactory().NewUserLoginRequest("", "")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		tokens, err := d.UserUsecase.Login(ctx, payload)
		if err != nil {
			return toHTTPError(logger, err)
		}

		// set refresh token as cookie
		cookie := new(http.Cookie)
		cookie.Name = refreshTokenCookie
		cookie.Value = tokens.RefreshToken
		if d.SecureRefreshToken {
			cookie.Secure = true
		}
		c.SetCookie(cookie)

		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"access_token": tokens.AccessToken}, "  ")
	}
}

func (d *UserDispatcher) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		cookie, err := c.Cookie(refreshTokenCookie)
		if err != nil {
			return toHTTPError(logger, err)
		}

		payload := dto.NewFactory().NewUserRefreshRequest(cookie.Value)
		tokens, err := d.UserUsecase.Refresh(ctx, payload)
		if err != nil {
			return toHTTPError(logger, err)
		}

		// set refresh token as cookie
		newCookie := new(http.Cookie)
		newCookie.Name = refreshTokenCookie
		newCookie.Value = tokens.RefreshToken
		if d.SecureRefreshToken {
			newCookie.Secure = true
		}
		c.SetCookie(newCookie)

		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"access_token": tokens.AccessToken}, "  ")
	}
}

func toHTTPError(logger *log.Logger, err error) error {
	// errors defined in usecase
	switch {
	case errors.Is(err, user.ErrInvalidRequest):
		return echo.NewHTTPError(http.StatusBadRequest)

	case errors.Is(err, user.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound)

	case errors.Is(err, user.ErrSystemError):
		logger.WithError(err).Error()
		return echo.NewHTTPError(http.StatusInternalServerError)

	case errors.Is(err, user.ErrDatabaseError):
		logger.WithError(err).Error()
		return echo.NewHTTPError(http.StatusInternalServerError)

	case errors.Is(err, user.ErrUnauthorized):
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	logger.WithError(err).Error()
	return echo.NewHTTPError(http.StatusInternalServerError)
}

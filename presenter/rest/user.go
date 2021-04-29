package rest

import (
	"errors"
	"net/http"

	"github.com/org39/webapp-tutorial-backend/presenter/rest/rr"
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

		payload := rr.NewFactory().NewUserSignUpRequest("", "")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user, tokens, err := d.UserUsecase.SignUp(ctx, payload.Email, payload.PlainPassword)
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

		return c.JSON(http.StatusCreated,
			rr.NewFactory().NewUserSignUpResponse(user.Email, user.CreatedAt, tokens.AccessToken),
		)
	}
}

func (d *UserDispatcher) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		logger := d.Logger.LoggerWithSpan(ctx)

		payload := rr.NewFactory().NewUserLoginRequest("", "")
		if err := c.Bind(payload); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		tokens, err := d.UserUsecase.Login(ctx, payload.Email, payload.PlainPassword)
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

		return c.JSON(http.StatusOK,
			rr.NewFactory().NewUserLoginResponse(tokens.AccessToken),
		)
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

		payload := rr.NewFactory().NewUserRefreshRequest(cookie.Value)
		tokens, err := d.UserUsecase.Refresh(ctx, payload.RefreshToken)
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

		return c.JSON(http.StatusOK,
			rr.NewFactory().NewUserRefreshResponse(tokens.AccessToken),
		)
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

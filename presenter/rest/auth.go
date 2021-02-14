package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/org39/webapp-tutorial-backend/entity/dto"
	"github.com/org39/webapp-tutorial-backend/usecase/auth"

	"github.com/labstack/echo/v4"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

const (
	authHeader = "Authorization"
	bearer     = "Bearer"
)

type AuthMiddleware struct {
	AuthUsercase auth.Usecase `inject:""`
	Logger       *log.Logger  `inject:""`
}

type AuthorizedContext struct {
	echo.Context
	id string
}

func (c *AuthorizedContext) ID() string {
	return c.id
}

func (a *AuthMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			logger := a.Logger.LoggerWithSpan(ctx)

			// Get access token from header
			authValue := req.Header.Get(authHeader)
			if len(authValue) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}
			token := extractBearerToken(authValue)
			if len(token) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "access token is required")
			}

			// verify token
			id, err := a.AuthUsercase.VerifyToken(ctx, dto.NewFactory().NewAuthVerifyRequest(token))
			switch {
			case errors.Is(err, auth.ErrUnauthorized):
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			case errors.Is(err, auth.ErrSystemError):
				logger.WithField("error", err).Error("")
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			case err != nil:
				logger.WithField("error", err).Error("")
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			// if token is valid, process request
			authorizedContext := &AuthorizedContext{c, id}
			err = next(authorizedContext)
			if err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}

func extractBearerToken(v string) string {
	invalidToken := ""
	parts := strings.Split(v, " ")

	switch {
	case len(parts) != 2:
		return invalidToken
	case parts[0] != bearer:
		return invalidToken
	}

	return parts[1]
}

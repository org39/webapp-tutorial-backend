package router

import (
	"strings"

	"github.com/HatsuneMiku3939/ocecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

func New() (*echo.Echo, error) {
	e := echo.New()

	e.Use(loggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(ocecho.OpenCensusMiddleware(
		ocecho.OpenCensusConfig{
			Skipper: func(c echo.Context) bool {
				// skip healthcheck endpoint
				return strings.HasPrefix(c.Path(), "/_health")
			},
			TraceOptions: ocecho.TraceOptions{
				IsPublicEndpoint: false,
				Propagation:      &b3.HTTPFormat{},
				StartOptions:     trace.StartOptions{},
			},
		},
	))

	e.HideBanner = true
	return e, nil
}

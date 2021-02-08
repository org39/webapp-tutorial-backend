package router

import (
	"time"

	"github.com/org39/webapp-tutorial-backend/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func loggerMiddleware(logger *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()
			ctx := req.Context()

			duration := time.Since(start)
			status := res.Status

			traceLogger := logger.LoggerWithSpan(ctx)
			accessLog := traceLogger.WithFields(logrus.Fields{
				"http_status":        status,
				"http_host":          req.Host,
				"http_latency":       float64(duration) / float64(1e6),
				"http_latency_human": duration.String(),
				"http_method":        req.Method,
				"http_uri":           req.RequestURI,
				"http_remote_ip":     c.RealIP(),
			})

			switch {
			case status >= 500:
				accessLog.Error()
			case status >= 400:
				accessLog.Warn()
			case status >= 300:
				accessLog.Debug()
			default:
				accessLog.Debug()
			}

			return nil
		}
	}
}

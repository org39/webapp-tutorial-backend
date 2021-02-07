package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	livenessEndpoint  = "/_healthy"
	readinessEndpoint = "/_healthz"
)

type Dispatchable interface {
	Dispatch(e *echo.Echo)
}

type Dispatcher struct {
	dispatchers []Dispatchable

	liveness  echo.HandlerFunc
	readiness echo.HandlerFunc
}

func NewDispatcher(options ...func(*Dispatcher) error) (*Dispatcher, error) {
	d := &Dispatcher{
		dispatchers: []Dispatchable{},
		liveness:    defaultHealthcheck(),
		readiness:   defaultHealthcheck(),
	}

	for _, option := range options {
		if err := option(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func WithLivenessCheck(h echo.HandlerFunc) func(*Dispatcher) error {
	return func(d *Dispatcher) error {
		d.liveness = h
		return nil
	}
}

func WithReadinessCheck(h echo.HandlerFunc) func(*Dispatcher) error {
	return func(d *Dispatcher) error {
		d.readiness = h
		return nil
	}
}

func (d *Dispatcher) Dispatch(e *echo.Echo) {
	// mount healthcheck endpoint
	e.GET(livenessEndpoint, d.liveness)
	e.GET(readinessEndpoint, d.readiness)

	for _, subr := range d.dispatchers {
		subr.Dispatch(e)
	}
}

func defaultHealthcheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

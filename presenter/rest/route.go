package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	healthcheckEndpoint = "/_health"
)

type Dispatchable interface {
	Dispatch(e *echo.Echo)
}

type Dispatcher struct {
	dispatchers []Dispatchable
	healthcheck echo.HandlerFunc
}

func NewDispatcher(options ...func(*Dispatcher) error) (*Dispatcher, error) {
	d := &Dispatcher{
		dispatchers: []Dispatchable{},
		healthcheck: nil,
	}

	for _, option := range options {
		if err := option(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func WithHealthcheck(h echo.HandlerFunc) func(*Dispatcher) error {
	return func(d *Dispatcher) error {
		d.healthcheck = h
		return nil
	}
}

func (d *Dispatcher) Dispatch(e *echo.Echo) {
	// mount healthcheck endpoint
	if d.healthcheck == nil {
		e.GET(healthcheckEndpoint, defaultHealthcheck())
	} else {
		e.GET(healthcheckEndpoint, d.healthcheck)
	}

	for _, subr := range d.dispatchers {
		subr.Dispatch(e)
	}
}

func defaultHealthcheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

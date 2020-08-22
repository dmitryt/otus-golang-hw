package app

import (
	"errors"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/service"
)

var ErrUnrecognizedServiceType = errors.New("cannot create service, because type was not recognized. Supported types: http, grpc")

type App struct {
	r repository.Base
	c *config.Config
}

func New(c *config.Config, r repository.Base) (*App, error) {
	return &App{c: c, r: r}, nil
}

func runService(addr string, stype service.TransportType, r repository.Base) error {
	s := service.New(stype, r)
	if s == nil {
		return ErrUnrecognizedServiceType
	}
	return s.Run(addr)
}

func (app *App) Run() error {
	errCh := make(chan error)

	go func() {
		if app.c.GRPCAddress != "" {
			errCh <- runService(app.c.GRPCAddress, service.GRPCType, app.r)
		}
	}()
	go func() {
		if app.c.HTTPAddress != "" {
			errCh <- runService(app.c.HTTPAddress, service.HTTPType, app.r)
		}
	}()
	return <-errCh
}

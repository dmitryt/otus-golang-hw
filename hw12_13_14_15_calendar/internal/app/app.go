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

func (app *App) Run() error {
	errCh := make(chan error)
	s := service.New(app.r)
	go func() {
		if app.c.GRPCAddress != "" {
			errCh <- s.Run(app.c.GRPCAddress)
		}
	}()
	go func() {
		if app.c.HTTPAddress != "" && app.c.GRPCAddress != "" {
			errCh <- s.HTTPProxy(app.c.GRPCAddress, app.c.HTTPAddress)
		}
	}()
	return <-errCh
}

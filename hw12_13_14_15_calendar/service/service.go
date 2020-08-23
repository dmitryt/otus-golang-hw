package service

import (
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	grpc "github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/service/server"
)

type Service interface {
	Run(string) error
	HTTPProxy(string, string) error
}

func New(r repository.Base) Service {
	return grpc.New(r)
}

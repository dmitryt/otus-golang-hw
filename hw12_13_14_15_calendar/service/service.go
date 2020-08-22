package service

import (
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	grpc "github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/service/grpc/server"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/service/http"
)

type Service interface {
	Run(string) error
}

type TransportType int

const (
	GRPCType TransportType = iota
	HTTPType
)

func New(stype TransportType, r repository.Base) Service {
	switch stype {
	case GRPCType:
		return grpc.New(r)
	case HTTPType:
		return http.New(r)
	default:
		return nil
	}
}

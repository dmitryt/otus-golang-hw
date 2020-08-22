package grpc

import (
	"context"
	"net"
	"time"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrConvertDBStructToProto = status.Error(codes.InvalidArgument, "cannot convert DB struct to proto one")
	ErrFetchingDataFromDB     = status.Error(codes.Internal, "error occurred during fetching data from DB")
	ErrUnsupportedRequest     = status.Error(codes.Unimplemented, "Request type is not supported")
)

type Service struct {
	r repository.Base
}

func New(r repository.Base) *Service {
	return &Service{r: r}
}

type QueryEventsType int

const (
	QueryDayEvents QueryEventsType = iota
	QueryWeekEvents
	QueryMonthEvents
)

func processEvents(repo repository.Base, queryType QueryEventsType, startTime time.Time) (result *EventsResponse, err error) {
	var dbEvents []repository.Event
	switch queryType {
	case QueryDayEvents:
		dbEvents, err = repo.GetDayEvents(startTime)
	case QueryWeekEvents:
		dbEvents, err = repo.GetWeekEvents(startTime)
	case QueryMonthEvents:
		dbEvents, err = repo.GetMonthEvents(startTime)
	default:
		return nil, ErrUnsupportedRequest
	}
	if err != nil {
		log.Error().Msgf("%s %s", ErrFetchingDataFromDB, err)
		return nil, ErrFetchingDataFromDB
	}
	for _, evt := range dbEvents {
		converted, err := ConvertEventToProto(evt)
		if err != nil {
			log.Error().Msgf("%s %s", ErrConvertDBStructToProto, err)
			return nil, ErrConvertDBStructToProto
		}
		result.Events = append(result.Events, converted)
	}
	return
}

func (s *Service) DayEvents(ctx context.Context, ts *timestamp.Timestamp) (result *EventsResponse, err error) {
	return processEvents(s.r, QueryDayEvents, time.Unix(ts.Seconds, int64(ts.Nanos)))
}

func (s *Service) WeekEvents(ctx context.Context, ts *timestamp.Timestamp) (result *EventsResponse, err error) {
	return processEvents(s.r, QueryWeekEvents, time.Unix(ts.Seconds, int64(ts.Nanos)))
}

func (s *Service) MonthEvents(ctx context.Context, ts *timestamp.Timestamp) (result *EventsResponse, err error) {
	return processEvents(s.r, QueryMonthEvents, time.Unix(ts.Seconds, int64(ts.Nanos)))
}

func (s *Service) CreateEvent(ctx context.Context, event *Event) (result *Event, err error) {
	evt, err := s.r.CreateEvent(*ConvertEventFromProto(event))
	if err != nil {
		log.Error().Msgf("[GRPC] Create Event: %s\n", err)
		return
	}
	log.Error().Msgf("[GRPC] Created Event: %+v\n", evt)
	result, err = ConvertEventToProto(evt)
	if err != nil {
		log.Error().Msgf("[GRPC] Create Event: %s\n", err)
	}
	return
}

func (s *Service) UpdateEvent(ctx context.Context, data *UpdateEventRequest) (result *Event, err error) {
	evt, err := s.r.UpdateEvent(data.Id, *ConvertEventFromProto(data.Event))
	if err != nil {
		return
	}
	return ConvertEventToProto(evt)
}

func (s *Service) DeleteEvent(ctx context.Context, data *DeleteEventRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.r.DeleteEvent(data.Id)
}

func (s *Service) Run(addr string) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()

	RegisterCalendarServer(server, s)

	log.Info().Msgf("Starting server on %s", lsn.Addr().String())

	return server.Serve(lsn)
}

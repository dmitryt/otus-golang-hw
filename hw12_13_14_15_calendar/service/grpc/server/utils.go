package grpc

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var ErrInvalidNullTime = errors.New("field value is invalid")

func ConvertTimeToTimestamp(ntime sql.NullTime) (result *timestamp.Timestamp, err error) {
	if !ntime.Valid {
		err = ErrInvalidNullTime
		return
	}
	value, err := ntime.Value()
	if err != nil {
		return nil, err
	}
	tvalue, ok := value.(time.Time)
	if !ok {
		err = ErrInvalidNullTime
		return
	}
	return ptypes.TimestampProto(tvalue)
}

// Didn't want to use reflection.
func ConvertEventToProto(evt repository.Event) (*Event, error) {
	result := &Event{
		ID:     evt.ID,
		UserID: evt.UserID,
		Title:  evt.Title,
	}
	if value, err := ptypes.TimestampProto(evt.StartDate); err == nil {
		result.StartDate = value
	} else {
		return nil, err
	}
	if value, err := ptypes.TimestampProto(evt.EndDate); err == nil {
		result.EndDate = value
	} else {
		return nil, err
	}
	if value, err := ConvertTimeToTimestamp(evt.NotifiedAt); err == nil {
		result.NotifiedAt = value
	} else {
		return nil, err
	}
	return result, nil
}

// Didn't want to use reflection.
func ConvertEventFromProto(evt *Event) *repository.Event {
	return &repository.Event{
		ID:        evt.ID,
		UserID:    evt.UserID,
		Title:     evt.Title,
		StartDate: time.Unix(evt.StartDate.GetSeconds(), int64(evt.StartDate.GetNanos())),
		EndDate:   time.Unix(evt.EndDate.GetSeconds(), int64(evt.EndDate.GetNanos())),
		NotifiedAt: sql.NullTime{
			Time:  time.Unix(evt.NotifiedAt.GetSeconds(), int64(evt.NotifiedAt.GetNanos())),
			Valid: true,
		},
	}
}

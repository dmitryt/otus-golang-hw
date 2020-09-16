package server

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var (
	ErrInvalidNullTime = errors.New("field value is invalid")
	ErrObjectIsNil     = errors.New("object is nil")
)

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
		ID:          int32(evt.ID),
		UserID:      int32(evt.UserID),
		Title:       evt.Title,
		Status:      evt.Status,
		NotifiedFor: int32(evt.NotifiedFor),
	}
	value, err := ptypes.TimestampProto(evt.StartDate)
	if err != nil {
		return nil, err
	}
	result.StartDate = value

	value, err = ptypes.TimestampProto(evt.EndDate)
	if err != nil {
		return nil, err
	}
	result.EndDate = value

	return result, nil
}

// Didn't want to use reflection.
func ConvertEventFromProto(evt *Event) (*repository.Event, error) {
	if evt == nil {
		return nil, ErrObjectIsNil
	}
	result := repository.Event{}
	if evt.ID != 0 {
		result.ID = int64(evt.ID)
	}
	if evt.UserID != 0 {
		result.UserID = int64(evt.UserID)
	}
	if evt.NotifiedFor != 0 {
		result.NotifiedFor = int64(evt.NotifiedFor)
	}
	if evt.Title != "" {
		result.Title = evt.Title
	}
	if evt.Status != "" {
		result.Status = evt.Status
	}
	if evt.StartDate != nil {
		result.StartDate = time.Unix(evt.StartDate.GetSeconds(), int64(evt.StartDate.GetNanos()))
	}
	if evt.EndDate != nil {
		result.EndDate = time.Unix(evt.EndDate.GetSeconds(), int64(evt.EndDate.GetNanos()))
	}
	return &result, nil
}

func ConvertUserToProto(data repository.User) (*User, error) {
	result := &User{
		ID:    int32(data.ID),
		Email: data.Email,
	}
	tmp, _ := data.FirstName.Value()
	if v, ok := tmp.(string); ok {
		result.FirstName = v
	}
	tmp, _ = data.LastName.Value()
	if v, ok := tmp.(string); ok {
		result.LastName = v
	}
	return result, nil
}

func ConvertUserFromProto(usr *User) (*repository.User, error) {
	if usr == nil {
		return nil, ErrObjectIsNil
	}
	result := repository.User{}
	if usr.ID != 0 {
		result.ID = int64(usr.ID)
	}
	if usr.Email != "" {
		result.Email = usr.Email
	}
	if usr.FirstName != "" {
		result.FirstName = sql.NullString{
			String: usr.FirstName,
			Valid:  true,
		}
	}
	if usr.LastName != "" {
		result.LastName = sql.NullString{
			String: usr.LastName,
			Valid:  true,
		}
	}
	return &result, nil
}

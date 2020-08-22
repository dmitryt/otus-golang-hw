package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
)

type Base interface {
	Connect(context.Context, *config.Config) error
	Close() error

	GetDayEvents(time.Time) ([]Event, error)
	GetWeekEvents(time.Time) ([]Event, error)
	GetMonthEvents(time.Time) ([]Event, error)

	CreateEvent(Event) (Event, error)
	UpdateEvent(int64, Event) (Event, error)
	DeleteEvent(int64) error
}

type Event struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	StartDate   time.Time      `db:"start_date"`
	StartTime   sql.NullTime   `db:"start_time"`
	EndDate     time.Time      `db:"end_date"`
	EndTime     sql.NullTime   `db:"end_time"`
	NotifiedAt  sql.NullTime   `db:"notified_at"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func New(repoType string) Base {
	switch repoType {
	case "psql":
		return NewPSQLRepo()
	case "memory":
		return NewMemoryRepo()
	}
	return nil
}
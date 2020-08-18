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

	GetEvents() ([]Event, error)
	GetEvent(int64) (Event, error)
	CreateEvent(Event) (Event, error)
	UpdateEvent(Event) (Event, error)
	DeleteEvent(int64) error
}

type Event struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	StartDate   string         `db:"start_date"`
	StartTime   time.Time      `db:"start_time"`
	EndDate     string         `db:"end_date"`
	EndTime     time.Time      `db:"end_time"`
	NotifiedAt  sql.NullString `db:"notified_at"`
	CreatedAt   string         `db:"created_at"`
	UpdatedAt   sql.NullString `db:"updated_at"`
}

func New(rtype string) Base {
	switch rtype {
	case "psql":
		return NewPSQLRepo()
	case "memory":
		return NewMemoryRepo()
	}
	return nil
}

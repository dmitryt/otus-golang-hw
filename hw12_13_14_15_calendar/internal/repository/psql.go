package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/jmoiron/sqlx"

	// is used for init postgres.
	_ "github.com/lib/pq"
)

var ErrDBOpen = errors.New("database open error")

type PSQLRepo struct {
	db            *sqlx.DB
	itemsPerQuery int
}

var insertQs = `INSERT INTO events
(user_id, title, description, start_date, start_time, end_date, end_time, notified_at)
VALUES
(:user_id, :title, :description, :start_date, :start_time, :end_date, :end_time, :notified_at)
RETURNING id`

var updateQs = `UPDATE events
	SET title=:title, description=:description, start_date=:start_date, start_time=:start_time, end_date=:end_date, end_time=:end_time, notified_at=:notified_at
	WHERE id = :id`

func getDSN(c *config.Config) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPass)
}

func (r *PSQLRepo) Connect(ctx context.Context, c *config.Config) (err error) {
	r.db, err = sqlx.Connect("postgres", getDSN(c))
	if err != nil {
		return fmt.Errorf("%s: %w", ErrDBOpen, err)
	}
	if c.DBMaxConn != 0 {
		r.db.SetMaxOpenConns(c.DBMaxConn)
	}
	if c.DBItemsPerQuery != 0 {
		r.itemsPerQuery = c.DBItemsPerQuery
	}
	return r.db.PingContext(ctx)
}

func (r *PSQLRepo) Close() error {
	return r.db.Close()
}

func NewPSQLRepo() *PSQLRepo {
	return &PSQLRepo{itemsPerQuery: 100}
}

func (r *PSQLRepo) GetEvents() (result []Event, err error) {
	err = r.db.Select(&result, "SELECT * FROM events ORDER BY start_date ASC LIMIT $1", r.itemsPerQuery)
	return
}

func (r *PSQLRepo) GetEvent(id int64) (result Event, err error) {
	err = r.db.Select(&result, "SELECT * FROM events where id = $1", id)
	return
}

func (r *PSQLRepo) DeleteEvent(id int64) (err error) {
	res, err := r.db.Exec("DELETE FROM events where id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrEventDelete, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrEventDelete, err)
	}
	if count == 0 {
		err = ErrEventDeleteFailed
	}
	return
}

func (r *PSQLRepo) CreateEvent(data Event) (event Event, err error) {
	res, err := r.db.NamedQuery(insertQs, data)
	if err != nil {
		return Event{}, fmt.Errorf("%s: %w", ErrEventCreate, err)
	}
	if res.Err() != nil {
		return Event{}, fmt.Errorf("%s: %w", ErrEventCreate, res.Err())
	}
	defer res.Rows.Close()

	event = Event{}
	for res.Rows.Next() {
		err := res.StructScan(&event)
		if err != nil {
			return Event{}, fmt.Errorf("%s: %w", ErrEventCreate, err)
		}
	}
	return event, nil
}

func (r *PSQLRepo) UpdateEvent(data Event) (event Event, err error) {
	res, err := r.db.NamedQuery(updateQs, data)
	if err != nil {
		return Event{}, fmt.Errorf("%s: %w", ErrEventUpdate, err)
	}
	if res.Err() != nil {
		return Event{}, fmt.Errorf("%s: %w", ErrEventUpdate, res.Err())
	}
	defer res.Rows.Close()

	event = Event{}
	for res.Rows.Next() {
		err := res.StructScan(&event)
		if err != nil {
			return Event{}, fmt.Errorf("%s: %w", ErrEventUpdate, err)
		}
	}
	return event, nil
}

package sqlstorage

import (
	"context"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Storage struct {
	Dsn string
	DB  *sqlx.DB
	Log *logger.Logger
}

func New(log *logger.Logger, conf config.StorageConf) *Storage {
	return &Storage{Log: log, Dsn: conf.Dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("postgres", s.Dsn)
	if err != nil {
		return err
	}

	s.DB = db
	return nil
}

func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		s.Log.Info(errors.Wrap(err, "err closing db connection").Error())
		return err
	}
	s.Log.Info("db connection gracefully closed")
	return nil
}

func (s *Storage) Add(ctx context.Context, event storage.Event) error {
	sql := `INSERT INTO events
(id, title, date, duration, description, userid, notify)
VALUES
    ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.DB.ExecContext(
		ctx,
		sql,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserId,
		event.Notify,
	)
	return err
}

func (s *Storage) Update(ctx context.Context, event storage.Event) error {
	sql := `UPDATE
    events
SET
    title = $2,
    date = $3,
    duration = $4,
    description = $5,
    userid = $6,
    notify = $7
WHERE
    id = $1`

	_, err := s.DB.ExecContext(
		ctx,
		sql,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserId,
		event.Notify,
	)
	return err
}

func (s *Storage) Del(ctx context.Context, event storage.Event) error {
	_, err := s.DB.ExecContext(ctx, "delete from events where id = &1", event.ID)
	return err
}

func (s *Storage) List(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	var ev []storage.Event
	err := s.DB.Select(&ev, "select * from events e where e.date = $1 and e.duration = $2", bg, bg.Sub(fn))
	return ev, err
}

func (s *Storage) ListByDay(ctx context.Context, dt time.Time) ([]storage.Event, error) {
	return s.List(ctx, dt, dt.AddDate(0, 0, 1))
}

func (s *Storage) ListByWeek(ctx context.Context, dt time.Time) ([]storage.Event, error) {
	return s.List(ctx, dt, dt.AddDate(0, 0, 7))
}

func (s *Storage) ListByMonth(ctx context.Context, dt time.Time) ([]storage.Event, error) {
	return s.List(ctx, dt, dt.AddDate(0, 1, 0))
}

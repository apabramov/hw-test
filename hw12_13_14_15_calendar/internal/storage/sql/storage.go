package sqlstorage

import (
	"context"
	_ "embed"
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
	Ctx context.Context
	Log *logger.Logger
}

var (
	//go:embed statements/insert.sql
	ins string

	//go:embed statements/update.sql
	upd string
)

func New(log *logger.Logger, conf config.StorageConf) *Storage {
	return &Storage{Log: log, Dsn: conf.Dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("postgres", s.Dsn)
	if err != nil {
		return err
	}

	s.DB = db
	s.Ctx = ctx
	return nil
}

func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		s.Log.Error(errors.Wrap(err, "err closing db connection").Error())
		return err
	}
	s.Log.Error("db connection gracefully closed")
	return nil
}

func (s *Storage) Add(event storage.Event) error {
	_, err := s.DB.ExecContext(
		s.Ctx,
		ins,
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

func (s *Storage) Upd(event storage.Event) error {
	_, err := s.DB.ExecContext(
		s.Ctx,
		upd,
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

func (s *Storage) Del(event storage.Event) error {
	_, err := s.DB.ExecContext(s.Ctx, "delete from events where id = &1", event.ID)
	return err
}

func (s *Storage) List(bg time.Time, fn time.Time) ([]storage.Event, error) {
	var ev []storage.Event
	err := s.DB.Select(&ev, "select * from events e where e.date = $1 and e.duration = $2", bg, bg.Sub(fn))
	return ev, err
}

func (s *Storage) ListByDay(dt time.Time) ([]storage.Event, error) {
	return s.List(dt, dt.AddDate(0, 0, 1))
}

func (s *Storage) ListByWeek(dt time.Time) ([]storage.Event, error) {
	return s.List(dt, dt.AddDate(0, 0, 7))
}

func (s *Storage) ListByMonth(dt time.Time) ([]storage.Event, error) {
	return s.List(dt, dt.AddDate(0, 1, 0))
}

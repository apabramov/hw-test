package sqlstorage

import (
	"context"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		event.Duration.Seconds(),
		event.Description,
		event.UserId,
		event.Notify.Seconds(),
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

func (s *Storage) Del(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "delete from events where id = &1", id)
	return err
}

func (s *Storage) Get(ctx context.Context, id string) (storage.Event, error) {
	var ev storage.EventPq
	err := s.DB.Get(&ev, "select * from events e where e.id = $1 ", id)
	d, err := ev.Duration.Duration()
	if err != nil {
		return storage.Event{}, err
	}
	n, err := ev.Notify.Duration()
	if err != nil {
		return storage.Event{}, err
	}
	return storage.Event{
		ID:          ev.ID,
		Title:       ev.Title,
		Date:        ev.Date,
		Duration:    d,
		Description: ev.Description,
		UserId:      ev.UserId,
		Notify:      n,
	}, err
}

func (s *Storage) List(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	var ev []storage.EventPq
	err := s.DB.Select(&ev, "select * from events e where e.date between $1 and $2", bg, fn)
	if err != nil {
		return nil, err
	}
	return convertEvent(ev)
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

func convertEvent(events []storage.EventPq) ([]storage.Event, error) {
	ev := make([]storage.Event, 0)

	for _, e := range events {
		d, err := e.Duration.Duration()
		if err != nil {
			return nil, err
		}
		n, err := e.Notify.Duration()
		if err != nil {
			return nil, err
		}
		ev = append(ev, storage.Event{
			ID:          e.ID,
			Title:       e.Title,
			Date:        e.Date,
			Duration:    d,
			Description: e.Description,
			UserId:      e.UserId,
			Notify:      n,
		})
	}
	return ev, nil
}

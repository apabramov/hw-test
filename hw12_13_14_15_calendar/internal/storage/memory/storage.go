package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[uuid.UUID]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[uuid.UUID]storage.Event)}
}

func (s *Storage) Add(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; ok {
		return storage.ErrExists
	}
	s.events[e.ID] = e
	return nil
}

func (s *Storage) Update(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}

func (s *Storage) Del(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return storage.ErrNotExists
	}
	delete(s.events, e.ID)
	return nil
}

func (s *Storage) List(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	ev := make([]storage.Event, 0, len(s.events))
	for _, e := range s.events {
		if e.Date.After(bg) && e.Date.Before(fn) {
			ev = append(ev, e)
		}
	}
	return ev, nil
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

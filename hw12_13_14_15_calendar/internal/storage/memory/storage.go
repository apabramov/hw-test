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

func (s *Storage) Connect(ctx context.Context) error {
	return nil
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

func (s *Storage) Del(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ud := uuid.MustParse(id)
	if _, ok := s.events[ud]; !ok {
		return storage.ErrNotExists
	}
	delete(s.events, ud)
	return nil
}

func (s *Storage) Get(ctx context.Context, id string) (storage.Event, error) {
	ud := uuid.MustParse(id)
	if _, ok := s.events[ud]; !ok {
		return storage.Event{}, storage.ErrNotExists
	}
	return s.events[ud], nil
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

func (s *Storage) ListByDay(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return s.List(ctx, bg, fn)
}

func (s *Storage) ListByWeek(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return s.List(ctx, bg, fn)
}

func (s *Storage) ListByMonth(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return s.List(ctx, bg, fn)
}

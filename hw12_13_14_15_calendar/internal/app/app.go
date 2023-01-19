package app

import (
	"context"
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type App struct {
	Log   Logger
	Store Storage
}

type Logger interface {
}

type Storage interface {
	Add(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, event storage.Event) error
	Del(ctx context.Context, event storage.Event) error
	ListByDay(ctx context.Context, dt time.Time) ([]storage.Event, error)
	ListByWeek(ctx context.Context, dt time.Time) ([]storage.Event, error)
	ListByMonth(ctx context.Context, dt time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{Log: logger, Store: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO

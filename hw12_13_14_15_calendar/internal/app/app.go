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
	Add(event storage.Event) error
	Upd(event storage.Event) error
	Del(event storage.Event) error
	ListByDay(dt time.Time) ([]storage.Event, error)
	ListByWeek(dt time.Time) ([]storage.Event, error)
	ListByMonth(dt time.Time) ([]storage.Event, error)
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

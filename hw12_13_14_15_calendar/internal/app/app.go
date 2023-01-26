package app

import (
	"context"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Log   Logger
	Store Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Storage interface {
	Add(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, event storage.Event) error
	Del(ctx context.Context, id string) error
	Get(ctx context.Context, is string) (storage.Event, error)
	ListByDay(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error)
	ListByWeek(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error)
	ListByMonth(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error)

	Connect(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{Log: logger, Store: storage}
}

func (a *App) AddEvent(ctx context.Context, e storage.Event) error {
	return a.Store.Add(ctx, e)
}

func (a *App) UpdateEvent(ctx context.Context, e storage.Event) error {
	return a.Store.Update(ctx, e)
}

func (a *App) DelEvent(ctx context.Context, id string) error {
	return a.Store.Del(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	return a.Store.Get(ctx, id)
}

func (a *App) ListByDayEvents(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return a.Store.ListByDay(ctx, bg, fn)
}

func (a *App) ListByWeekEvents(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return a.Store.ListByWeek(ctx, bg, fn)
}

func (a *App) ListByMonthEvents(ctx context.Context, bg time.Time, fn time.Time) ([]storage.Event, error) {
	return a.Store.ListByMonth(ctx, bg, fn)
}

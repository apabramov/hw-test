package app

import (
	"context"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	Log   Logger
	Store Storage
}

func NewScheduler(logger Logger, storage Storage) *Scheduler {
	return &Scheduler{
		Log:   logger,
		Store: storage,
	}
}

func (a *Scheduler) ListNotify(ctx context.Context) ([]storage.Event, error) {
	return a.Store.ListNotify(ctx, time.Now())
}

func (a *Scheduler) DeleteOutDate(ctx context.Context) error {
	return a.Store.DeleteOutDate(ctx, time.Now().Add(-time.Second*60))
}

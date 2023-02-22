package app

import (
	"context"
	"fmt"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/queue"
)

type Sender struct {
	Log   Logger
	Store Storage
}

func NewSender(logger Logger, storage Storage) *Sender {
	return &Sender{
		Log:   logger,
		Store: storage,
	}
}

func (a *Sender) SendNotify(ctx context.Context, n queue.Notification) error {
	e, err := a.Store.Get(ctx, n.GetId())
	if err != nil {
		return err
	}
	if !e.Sent {
		e.Sent = true
		a.Store.Update(ctx, e)
		a.Log.Info(fmt.Sprintf("ID: %v Title: %v Date: %v, UserId: %v", n.GetId(), n.GetTitle(), n.GetDate(), n.GetUserId()))
	}
	return nil
}

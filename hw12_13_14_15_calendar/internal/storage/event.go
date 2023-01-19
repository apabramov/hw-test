package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExists    = errors.New("event exists")
	ErrNotExists = errors.New("event not exists")
	ErrDateBusy  = errors.New("date busy")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Date        time.Time
	Duration    time.Duration
	Description string
	UserId      uuid.UUID
	Notify      time.Duration
}

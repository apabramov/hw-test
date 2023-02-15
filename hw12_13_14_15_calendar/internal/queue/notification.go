package queue

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserId uuid.UUID `json:"userid"`
}

func (n *Notification) GetId() string {
	return n.ID.String()
}

func (n *Notification) GetTitle() string {
	return n.Title
}

func (n *Notification) GetDate() string {
	return n.Date.String()
}

func (n *Notification) GetUserId() string {
	return n.UserId.String()
}

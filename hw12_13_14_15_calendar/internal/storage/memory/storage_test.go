package memorystorage

import (
	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	s := New()
	err := s.Add(e)
	require.NoError(t, err)

	t.Run("update", func(t *testing.T) {
		e.Description = "update description"
		err := s.Upd(e)
		require.NoError(t, err)
	})

	t.Run("count event", func(t *testing.T) {
		l, err := s.ListByDay(e.Date.Add(-10 * time.Second))
		require.NoError(t, err)
		require.Len(t, l, 1)
	})

	t.Run("list week", func(t *testing.T) {
		ew := storage.Event{
			ID:   uuid.New(),
			Date: time.Now().AddDate(0, 0, 1),
		}
		s.Add(ew)
		l, err := s.ListByWeek(e.Date.Add(-10 * time.Second))
		require.NoError(t, err)
		require.Len(t, l, 2)
	})

	t.Run("list month", func(t *testing.T) {
		l, err := s.ListByMonth(e.Date.Add(-10 * time.Second))
		require.NoError(t, err)
		require.Len(t, l, 2)
	})

	t.Run("delete", func(t *testing.T) {
		err := s.Del(e)
		require.NoError(t, err)
	})
}

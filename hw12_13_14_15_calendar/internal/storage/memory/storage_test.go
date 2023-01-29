package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAddStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)
}

func TestUpdateStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)

	t.Run("update", func(t *testing.T) {
		e.Description = "update description"
		err := s.Update(ctx, e)
		require.NoError(t, err)

		l, err := s.ListByDay(ctx, e.Date.Add(-10*time.Second), e.Date.Add(-10*time.Second).AddDate(0, 0, 1))
		require.NoError(t, err)
		require.True(t, l[0].Description == "update description")
	})
}

func TestCountStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)
	t.Run("count event", func(t *testing.T) {
		l, err := s.ListByDay(ctx, e.Date.Add(-10*time.Second), e.Date.Add(-10*time.Second).AddDate(0, 0, 1))
		require.NoError(t, err)
		require.Len(t, l, 1)
	})
}

func TestListWeekStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)

	t.Run("list week", func(t *testing.T) {
		ew := storage.Event{
			ID:   uuid.New(),
			Date: time.Now().AddDate(0, 0, 1),
		}
		s.Add(ctx, ew)
		l, err := s.ListByWeek(ctx, e.Date.Add(-10*time.Second), e.Date.Add(-10*time.Second).AddDate(0, 0, 7))
		require.NoError(t, err)
		require.Len(t, l, 2)
	})
}

func TestListMonthStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)
	t.Run("list month", func(t *testing.T) {
		l, err := s.ListByMonth(ctx, e.Date.Add(-10*time.Second), e.Date.Add(-10*time.Second).AddDate(0, 1, 0))
		require.NoError(t, err)
		require.Len(t, l, 1)
	})
}

func TestDeleteStorage(t *testing.T) {
	e := storage.Event{
		ID:          uuid.New(),
		Title:       "Testing Title",
		Date:        time.Now(),
		Duration:    60 * time.Minute,
		Description: "Description",
		UserId:      uuid.New(),
		Notify:      5 * time.Minute,
	}

	ctx := context.Background()

	s := New()
	err := s.Add(ctx, e)
	require.NoError(t, err)
	t.Run("delete", func(t *testing.T) {
		err := s.Del(ctx, e.ID.String())
		require.NoError(t, err)
	})
}

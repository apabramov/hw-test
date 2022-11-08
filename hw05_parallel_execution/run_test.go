package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks eventually", func(t *testing.T) {
		var runTasksCount, finishTasksCount int32
		workersCount := 5
		tasks := make([]Task, 0, workersCount)
		waitCh := make(chan struct{})

		for i := 0; i < workersCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				<-waitCh
				atomic.AddInt32(&finishTasksCount, 1)
				return nil
			})
		}

		cntErr := make(chan error)

		go func() {
			cntErr <- Run(tasks, workersCount, workersCount)
		}()

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&runTasksCount) == int32(workersCount)
		}, time.Second, time.Millisecond)

		close(waitCh)

		require.NoError(t, <-cntErr)
	})
}

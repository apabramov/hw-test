package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var cntErr int32
	taskIn := make(chan Task)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go do(&wg, taskIn, &cntErr)
	}

	for _, t := range tasks {
		if atomic.LoadInt32(&cntErr) >= int32(m) {
			break
		}
		taskIn <- t
	}

	close(taskIn)
	wg.Wait()

	if cntErr >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func do(wg *sync.WaitGroup, taskIn chan Task, cnt *int32) {
	defer wg.Done()

	for t := range taskIn {
		if err := t(); err != nil {
			atomic.AddInt32(cnt, 1)
		}
	}
}

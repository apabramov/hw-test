package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Worker struct {
	in     chan Task
	outErr chan error
	done   chan struct{}
	wg     *sync.WaitGroup
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	w := Worker{
		in:     make(chan Task),
		outErr: make(chan error),
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
	}

	go func(w Worker) {
		defer close(w.in)

		for i := range tasks {
			if IsClosed(w.done) {
				return
			}
			w.in <- tasks[i]
		}
	}(w)

	go func(w Worker) {
		for i := 0; i < n; i++ {
			w.wg.Add(1)
			go do(w)
		}
		w.wg.Wait()
		close(w.outErr)
	}(w)

	cnt := 0
	err := false
	for res := range w.outErr {
		cnt++
		fmt.Println(res)
		if cnt == m {
			close(w.done)
			err = true
		}
	}

	if err {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func IsClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func do(w Worker) {
	defer w.wg.Done()

	for t := range w.in {
		if IsClosed(w.done) {
			return
		}
		if err := t(); err != nil {
			w.outErr <- err
		}
	}
}

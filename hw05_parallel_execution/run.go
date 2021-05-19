package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wg *sync.WaitGroup, in chan Task, errCounter *int32) {
	defer wg.Done()
	for task := range in {
		err := task()
		if err != nil {
			atomic.AddInt32(errCounter, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(taskList []Task, n, m int) error {
	if m < 1 {
		m = 1
	}
	if n < 1 {
		n = 1
	}
	wg := &sync.WaitGroup{}
	in := make(chan Task)
	var errCounter int32
	wg.Add(n)

	for i := 0; i < n; i++ {
		go worker(wg, in, &errCounter)
	}
	var result error
	for _, task := range taskList {
		if int(atomic.LoadInt32(&errCounter)) >= m {
			break
		}
		in <- task
	}
	close(in)
	wg.Wait()
	if int(errCounter) >= m {
		result = ErrErrorsLimitExceeded
	}
	return result
}

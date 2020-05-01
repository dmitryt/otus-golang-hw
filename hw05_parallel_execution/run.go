package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorsCounter struct {
	value int
	mu sync.Mutex
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	var wg sync.WaitGroup

	queueCh := make(chan int, n)
	quitCh := make(chan int)
	errorsCnt := errorsCounter{}

	for _, task := range tasks {
		select {
		case <- quitCh:
			break
		case queueCh <- 1:
			wg.Add(1)
			go func(task Task){
				defer wg.Done()
				releaseQueue := true
				if (task() != nil) {
					errorsCnt.mu.Lock()
					errorsCnt.value++
					if m == errorsCnt.value - 1 {
						releaseQueue = false
						close(quitCh)
					}
					errorsCnt.mu.Unlock()
				}
				if releaseQueue {
					<- queueCh
				}
			}(task)
		}
	}

	// Wait, when all goroutines finish
	wg.Wait()

	if errorsCnt.value > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

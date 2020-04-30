package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func prepareClosedChannel(cap int) chan int {
	channel := make(chan int, cap)
	for i := 0; i < cap; i++ {
		channel <- 1
	}
	close(channel)
	return channel
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	var result error = nil
	var wg sync.WaitGroup

	queueCh := make(chan int, n)
	quitCh := make(chan error)
	errorCh := prepareClosedChannel(m)

	for _, task := range tasks {
		select {
		case result = <-quitCh:
			break
		case queueCh <- 1:
			wg.Add(1)
			go func() {
				result := task()
				wg.Done()
				if result != nil {
					if _, ok := <-errorCh; !ok {
						quitCh <- ErrErrorsLimitExceeded
					}
				}
				<-queueCh
			}()
		}
	}
	// Waiting, when all tasks in queue are finished
	wg.Wait()
	return result
}

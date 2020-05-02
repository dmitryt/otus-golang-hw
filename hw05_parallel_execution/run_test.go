package hw05_parallel_execution //nolint:golint,stylecheck

import (
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
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("should fail on first error, when amount of allowed errors is 0", func(t *testing.T) {
		tasksCount := 10
		errTaskIndex := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			tasks = append(tasks, func (i int) (func() error) {
				return func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					if i == errTaskIndex {
						return fmt.Errorf("error from task %d", i)
					}
					return nil
				}}(i))
		}

		workersCount := 3
		maxErrorsCount := 0

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t,
			int32(workersCount+maxErrorsCount), runTasksCount, "extra tasks were started")
	})

	t.Run("should handle tasks with different duration correctly", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < 2; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(100))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(10))
			atomic.AddInt32(&runTasksCount, 1)
			return fmt.Errorf("error from task")
		})

		for i := 0; i < 7; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(10))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 3
		maxErrorsCount := 1

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, result)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("should handle errors correctly, when last task exceeds errors' limit", func(t *testing.T) {
		successTasksCount := 10
		errorTasksCount := 2
		tasks := make([]Task, 0, successTasksCount)

		var runTasksCount int32

		for i := 0; i < successTasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		for i := 0; i < errorTasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 3
		maxErrorsCount := 1

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
	})
}

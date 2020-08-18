package repository

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func wrap(vs ...interface{}) []interface{} {
	return vs
}

func TestMemoryRepo(t *testing.T) {
	t.Run("empty repo", func(t *testing.T) {
		r := NewMemoryRepo()

		_, err := r.GetEvent(1)
		require.EqualError(t, err, fmt.Sprintf("%s", ErrEventNotFound))

		_, err = r.GetEvent(2)
		require.EqualError(t, err, fmt.Sprintf("%s", ErrEventNotFound))
	})

	t.Run("create event", func(t *testing.T) {
		r := NewMemoryRepo()

		startDate := time.Now().Format("2006-01-02")
		endDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
		event, err := r.CreateEvent(Event{UserID: 1, StartDate: startDate, EndDate: endDate})
		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, event, r.storage[event.ID])
	})

	t.Run("update event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := time.Now().Format("2006-01-02")
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		endDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
		event, err := r.UpdateEvent(Event{ID: id, UserID: 1, EndDate: endDate})

		require.Nil(t, err)
		require.Equal(t, event, r.storage[id])
	})

	t.Run("remove event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := time.Now().Format("2006-01-02")
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		err := r.DeleteEvent(id)
		require.Nil(t, err)

		_, ok := r.storage[id]
		require.False(t, ok)
	})

	t.Run("remove unavailable event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := time.Now().Format("2006-01-02")
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		err := r.DeleteEvent(111)
		require.EqualError(t, err, fmt.Sprintf("%s", ErrEventNotFound))
	})
}

func TestRepoMultithreading(t *testing.T) {
	t.Run("test multithreading", func(t *testing.T) {
		r := NewMemoryRepo()
		wg := &sync.WaitGroup{}
		wg.Add(2)

		startDate := time.Now().Format("2006-01-02")
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				r.CreateEvent(Event{UserID: 1, StartDate: startDate})
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				r.DeleteEvent(int64(i))
			}
		}()

		wg.Wait()
	})
}

package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func wrap(vs ...interface{}) []interface{} {
	return vs
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(10)

		require.False(t, c.Set("aaa", 100))
		require.False(t, c.Set("bbb", 200))
		require.False(t, c.Set("ccc", 300))

		// Check values in the cache
		require.Equal(t, []interface{}{100, true}, wrap(c.Get("aaa")))
		require.Equal(t, []interface{}{200, true}, wrap(c.Get("bbb")))
		require.Equal(t, []interface{}{300, true}, wrap(c.Get("ccc")))

		c.Clear()

		// Check values in the cache
		require.Equal(t, []interface{}{nil, false}, wrap(c.Get("aaa")))
		require.Equal(t, []interface{}{nil, false}, wrap(c.Get("bbb")))
		require.Equal(t, []interface{}{nil, false}, wrap(c.Get("ccc")))
	})

	t.Run("check cache capacity", func(t *testing.T) {
		c := NewCache(4)

		require.False(t, c.Set("aaa", 100))
		require.False(t, c.Set("bbb", 200))
		require.False(t, c.Set("ccc", 300))
		require.False(t, c.Set("ddd", 400))
		require.False(t, c.Set("eee", 500))

		// Item should be removed from the cache
		require.Equal(t, []interface{}{nil, false}, wrap(c.Get("aaa")))
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

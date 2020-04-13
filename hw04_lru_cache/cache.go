package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mux      sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity, queue: NewList(), items: make(map[Key]*listItem)}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	item, found := c.items[key]
	// it's strange, that I cannot return c.items[key] directly
	if !found {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, true
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	item, found := c.items[key]
	if found {
		// Looks like I cannot reassign value in a such way
		// item.Value.(cacheItem).value = value
		c.queue.Remove(item)
	}
	c.items[key] = c.queue.PushFront(cacheItem{key: key, value: value})
	if c.queue.Len() > c.capacity {
		itemToRemove := c.queue.Back()
		delete(c.items, itemToRemove.Value.(cacheItem).key)
		c.queue.Remove(itemToRemove)
	}
	return found
}

func (c *lruCache) Clear() {
	for key, item := range c.items {
		c.queue.Remove(item)
		delete(c.items, key)
	}
}

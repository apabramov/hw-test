package hw04lrucache

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
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	ch := &cacheItem{key: key, value: value}
	c.mu.Lock()

	defer c.mu.Unlock()
	if _, ok := c.items[key]; ok {
		c.items[key] = &ListItem{Value: value}
		c.queue.PushFront(ch)
		return true
	}
	c.items[key] = &ListItem{Value: value}
	c.queue.PushFront(ch)
	if c.queue.Len() > c.capacity {
		b := c.queue.Back()
		c.queue.Remove(b)
		delete(c.items, b.Value.(*cacheItem).key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()

	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		ch := &cacheItem{key: key, value: val}
		b := &ListItem{Value: ch}
		c.queue.MoveToFront(b)
		return val.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = nil
	c.items = nil
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

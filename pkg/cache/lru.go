package cache

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

type lruCache struct {
	mu  *sync.Mutex
	lru *lru.Cache
}

func newLRUCache(maxEntries int) *lruCache {
	return &lruCache{
		mu:  new(sync.Mutex),
		lru: lru.New(maxEntries),
	}
}

func (c *lruCache) Get(key int64) (value interface{}, exists bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lru.Get(key)
}

func (c *lruCache) Remove(key int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Remove(key)
}

func (c *lruCache) Add(key int64, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Add(key, value)
}

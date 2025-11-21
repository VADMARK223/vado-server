package app

import (
	"sync"
	"time"
)

type item struct {
	value any
	exp   time.Time
}

type LocalCache struct {
	mu    sync.RWMutex
	items map[any]item
}

func NewLocalCache() *LocalCache {
	return &LocalCache{
		items: make(map[any]item),
	}
}

func (c *LocalCache) Set(key, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item{
		value: value,
		exp:   time.Now().Add(ttl),
	}
}

func (c *LocalCache) Get(key any) (any, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(it.exp) {
		return nil, false
	}
	return it.value, true
}

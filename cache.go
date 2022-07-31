package GOcache

import (
	"GOcache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.LruCache
	cachebytes int64
}

func (c *cache) add(key string, value Byteview) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cachebytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value Byteview, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if ele, ok := c.lru.Get(key); ok {
		return ele.(Byteview), ok
	}
	return
}

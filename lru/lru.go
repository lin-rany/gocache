package lru

import "container/list"

type LruCache struct {
	maxBytes  int64
	nBytes    int64
	data      *list.List
	cachemap  map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type Entry struct {
	key   string
	value Value
}

// New cache
func New(maxbytes int64, onEvicted func(key string, value Value)) *LruCache {
	return &LruCache{
		maxBytes:  maxbytes,
		nBytes:    0,
		data:      list.New(),
		cachemap:  make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get data from cache
func (c *LruCache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cachemap[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		return kv.value, true
	}
	return
}

// remove oldest data
func (c *LruCache) RemoveOldest() {
	ele := c.data.Back()
	if ele != nil {
		kv := ele.Value.(*Entry)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		delete(c.cachemap, ele.Value.(*Entry).key)
		c.data.Remove(ele)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// add data to cache
func (c *LruCache) Add(key string, value Value) {
	if ele, ok := c.cachemap[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.data.PushFront(&Entry{key: key, value: value})
		c.cachemap[key] = ele
		c.nBytes += int64(value.Len()) + int64(len(key))
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *LruCache) Len() int {
	return c.data.Len()
}

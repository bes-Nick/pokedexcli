package pokecache

import (
	"sync"
	"time"
)

// storing entire cache in memory
type Cache struct {
	cache map[string]cacheEntry
	mux   *sync.Mutex
}

// cacheEntry
type cacheEntry struct {
	val      []byte
	createAt time.Time
}

// helper method to create a cache
func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

//a add method to cache

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cacheEntry{
		val:      val,
		createAt: time.Now().UTC(),
	}
}

// get method

func (c *Cache) Get(key string) ([]byte, bool) {
	cacheE, ok := c.cache[key]
	return cacheE.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

// method for delete from cache
func (c *Cache) reap(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	timeAgo := time.Now().UTC().Add(-interval)
	for k, v := range c.cache {
		if v.createAt.Before(timeAgo) {
			delete(c.cache, k)
		}
	}
}

package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	Entries map[string]cacheEntry
	mu      *sync.RWMutex
}

// Create a new cache
func NewCache(interval time.Duration) Cache {
	cache := Cache{
		Entries: make(map[string]cacheEntry),
		mu:      &sync.RWMutex{},
	}
	go cache.reapLoop(interval)

	return cache

}

func (c *Cache) Add(key string, value *[]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		value:     *value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.Entries[key]
	return entry.value, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, entry := range c.Entries {
		if entry.createdAt.Before(now.Add(-last)) {
			delete(c.Entries, name)
		}
	}
}

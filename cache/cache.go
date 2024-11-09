package cache

import (
	"sync"
	"time"
)

// Entry represents a single entry in the cache
type Entry struct {
	Value      interface{}
	Expiration int64
}

// Cache is a basic in-memory cache
type Cache struct {
	mu    sync.RWMutex
	store map[string]Entry
}

// NewCache creates a new instance of Cache
func NewCache() *Cache {
	return &Cache{
		store: make(map[string]Entry),
	}
}

// Set adds an entry to the cache with a specific expiration time
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = Entry{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Get retrieves an entry from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.store[key]
	if !found || time.Now().UnixNano() > entry.Expiration {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes an entry from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Flush clears the entire cache
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]Entry)
}

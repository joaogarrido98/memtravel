package cache

import (
	"sync"
	"time"
)

// CacheEntry represents a single entry in the cache
type CacheEntry struct {
	Value      interface{}
	Expiration int64
}

// SimpleCache is a basic in-memory cache
type SimpleCache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

// NewSimpleCache creates a new instance of SimpleCache
func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		store: make(map[string]CacheEntry),
	}
}

// Set adds an entry to the cache with a specific expiration time
func (c *SimpleCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = CacheEntry{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Get retrieves an entry from the cache
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.store[key]
	if !found || time.Now().UnixNano() > entry.Expiration {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes an entry from the cache
func (c *SimpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Flush clears the entire cache
func (c *SimpleCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]CacheEntry)
}

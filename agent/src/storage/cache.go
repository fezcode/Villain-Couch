package storage

import (
	"sync"
)

// Cache is a generic, thread-safe in-memory key-value store.
// T can be any type.
type Cache[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// NewCache creates and returns a new generic Cache instance.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		items: make(map[string]T),
	}
}

// Set adds or updates an item in the cache.
func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// Get retrieves an item from the cache.
// It returns the value and a boolean indicating if the key was found.
// If the key is not found, the zero value of type T is returned.
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	return item, found
}

// Delete removes an item from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Keys returns a slice of all keys in the cache.
// This operation is thread-safe.
func (c *Cache[T]) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Pre-allocate a slice with the right capacity for efficiency.
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

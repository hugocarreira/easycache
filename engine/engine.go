package engine

import "time"

// Engine defines the core behavior of a cache system.
//
// This interface abstracts different caching strategies, including FIFO, LRU, LFU, and TTL-based caches.
// Implementations of this interface determine how items are stored, retrieved, and evicted.
type Engine interface {
	// Get retrieves a value from the cache by its key.
	// Returns (value, true) if the key exists, otherwise returns (nil, false).
	Get(key string) (any, bool)

	// Set stores a key-value pair in the cache.
	// If the key already exists, its value is updated.
	Set(key string, value any)

	// SetWithTTL stores a key-value pair in the cache with an expiration time.
	// This method is only relevant for TTL-based caches.
	SetWithTTL(key string, value any, expiresAt time.Time)

	// Delete removes a key-value pair from the cache.
	Delete(key string)

	// Has checks whether a given key exists in the cache.
	// Returns true if the key is present and has not expired (for TTL-based caches).
	Has(key string) bool

	// Len returns the number of items currently stored in the cache.
	Len() int

	// IsExpirable returns true if the cache supports TTL-based expiration.
	IsExpirable() bool

	// IsExpired checks whether a specific key has expired.
	IsExpired(key string) bool

	// Evict removes an item from the cache based on the eviction policy (FIFO, LRU, LFU).
	Evict()
}

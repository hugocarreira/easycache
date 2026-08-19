package cache

import (
	"runtime"
	"sync"
	"time"

	"github.com/hugocarreira/easycache/basic"
	"github.com/hugocarreira/easycache/fifo"
	"github.com/hugocarreira/easycache/lfu"
	"github.com/hugocarreira/easycache/lru"

	"github.com/hugocarreira/easycache/engine"
)

// EvictionPolicy defines the possible cache eviction strategies.
//
// The eviction policy determines how items are removed when the cache reaches
// its maximum size. The available policies are:
//
//   - Basic: No automatic eviction, items are removed only when they expire (TTL-based).
//   - FIFO: First-In, First-Out eviction; the oldest item is removed first.
//   - LRU: Least Recently Used eviction; the least accessed item is removed first.
//   - LFU: Least Frequently Used eviction; the item with the fewest accesses is removed first.
type EvictionPolicy int

const (
	Basic EvictionPolicy = iota
	FIFO
	LRU
	LFU
)

// Cache is the main structure that manages an in-memory key-value store
// with different eviction policies and optional TTL-based expiration.
//
// It acts as a wrapper around specific caching strategies such as FIFO, LRU, LFU,
// or a simple TTL-based cache. The eviction policy is defined in the CacheConfig.
//
// The Cache structure provides thread-safe access through its selected engine
// and includes built-in metrics for monitoring performance.
type Cache struct {
	// engine represents the selected cache strategy (FIFO, LRU, LFU, or Basic).
	// It implements the CacheInterface to allow dynamic eviction policies.
	engine engine.Engine

	// Config holds the configuration settings, such as eviction policy,
	// max size, and TTL (if applicable).
	config *Config

	// metrics tracks cache statistics, including hits and misses.
	metrics *Metrics

	stop      chan struct{}
	closeOnce sync.Once
}

func New(cfg *Config) *Cache {
	config := *defaultConfig()
	if cfg != nil {
		config = *cfg
	}

	if config.CleanupInterval <= 0 {
		config.CleanupInterval = 10 * time.Second
	}
	if config.MaxSize < 0 {
		config.MaxSize = 0
	}
	if config.EvictionPolicy < Basic || config.EvictionPolicy > LFU {
		config.EvictionPolicy = Basic
	}

	c := &Cache{
		config:  &config,
		metrics: NewMetrics(),
		stop:    make(chan struct{}),
	}

	switch config.EvictionPolicy {
	case LRU:
		c.engine = lru.New(config.MaxSize)
	case FIFO:
		c.engine = fifo.New(config.MaxSize)
	case LFU:
		c.engine = lfu.New(config.MaxSize)
	default:
		c.engine = basic.New(config.MaxSize, config.TTL, config.CleanupInterval)
	}

	if config.MemoryLimits > 0 && config.MemoryCheckInterval > 0 {
		go c.startCheckMemoryUsage()
	}

	return c
}

// startCheckMemoryUsage periodically monitors process heap allocation.
//
// If memory limits are set in CacheConfig, this function runs at the configured
// interval (`MemoryCheckInterval`). When memory usage exceeds `MemoryLimits`,
// the cache triggers cleanup to free up space
func (c *Cache) startCheckMemoryUsage() {
	if c.config.MemoryLimits == 0 {
		return
	}

	if c.config.MemoryCheckInterval <= 0 {
		return
	}

	ticker := time.NewTicker(c.config.MemoryCheckInterval)
	defer ticker.Stop()

	maxMem := c.config.MemoryLimits

	for {
		select {
		case <-ticker.C:
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			if mem.Alloc > maxMem {
				c.engine.Evict()
			}
		case <-c.stop:
			return
		}
	}
}

// Get retrieves a value from the cache by its key.
//
// If the key exists and has not expired, the function returns the value and true.
// If the key does not exist or has expired (in case of TTL-based eviction),
// the function returns nil and false. Additionally, cache hit/miss metrics
// are updated accordingly.
func (c *Cache) Get(key string) (any, bool) {
	elem, exists := c.engine.Get(key)

	if !exists {
		if c.config.Metrics {
			c.metrics.IncrementMisses()
		}
		return nil, false
	}

	if c.config.Metrics {
		c.metrics.IncrementHits()
	}

	return elem, true
}

// Set stores a key-value pair in the cache.
//
// If the key already exists, its value is updated. For FIFO, LRU, and LFU
// policies, a positive MaxSize is enforced atomically by the selected engine.
// Basic caches ignore MaxSize and use TTL expiration only. If TTL is positive,
// the item expires after the configured duration.
func (c *Cache) Set(key string, value string) {
	if c.engine.IsExpirable() {
		expiresAt := time.Time{}
		if c.config.TTL > 0 {
			expiresAt = time.Now().Add(c.config.TTL)
		}
		c.engine.SetWithTTL(key, value, expiresAt)
		return
	}

	c.engine.Set(key, value)
}

// Delete removes a key-value pair from the cache.
//
// If the key exists, it is removed from both the primary storage and any
// auxiliary structures (e.g., linked lists for LRU/FIFO or heaps for LFU).
// If the key does not exist, the function does nothing.
func (c *Cache) Delete(key string) {
	c.engine.Delete(key)
}

// Has checks whether a given key exists in the cache.
//
// Returns true if the key is present and has not expired (for TTL-based caches).
// If the key does not exist or has expired, it returns false.
func (c *Cache) Has(key string) bool {
	return c.engine.Has(key)
}

// Len returns the number of items currently stored in the cache.
//
// For TTL-based caches, only non-expired items are counted. In other eviction
// policies (FIFO, LRU, LFU), it returns the total number of stored items.
func (c *Cache) Len() int {
	return c.engine.Len()
}

func (c *Cache) Evict() {
	c.engine.Evict()
}

// Metrics returns a pointer to the cache's metrics instance.
//
// The metrics track cache performance, including hits and misses.
// The returned object is always non-nil; counters remain unchanged when metrics
// collection is disabled in the configuration.
func (c *Cache) Metrics() *Metrics {
	return c.metrics
}

// Close stops background cleanup and memory-monitoring goroutines.
// It is safe to call Close more than once.
func (c *Cache) Close() {
	c.closeOnce.Do(func() {
		close(c.stop)
		if closeable, ok := c.engine.(interface{ Close() }); ok {
			closeable.Close()
		}
	})
}

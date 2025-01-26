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
// The Cache structure provides thread-safe access with read/write locks and
// includes built-in metrics for monitoring performance.
type Cache struct {
	// lock ensures thread-safe access to the cache data.
	lock sync.RWMutex

	// engine represents the selected cache strategy (FIFO, LRU, LFU, or Basic).
	// It implements the CacheInterface to allow dynamic eviction policies.
	engine engine.Engine

	// Config holds the configuration settings, such as eviction policy,
	// max size, and TTL (if applicable).
	config *Config

	// metrics tracks cache statistics, including hits and misses.
	metrics *Metrics
}

func New(cfg *Config) *Cache {
	if cfg == nil {
		cfg = defaultConfig()
	}

	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = 10 * time.Second
	}

	c := &Cache{
		config:  cfg,
		metrics: NewMetrics(),
	}

	switch cfg.EvictionPolicy {
	case LRU:
		c.engine = lru.New(cfg.MaxSize)
	case FIFO:
		c.engine = fifo.New(cfg.MaxSize)
	case LFU:
		c.engine = lfu.New(cfg.MaxSize)
	default:
		c.engine = basic.New(cfg.MaxSize, cfg.TTL, cfg.CleanupInterval)
	}

	go c.startCheckMemoryUsage()

	return c
}

// startCheckMemoryUsage periodically monitors the cache's memory usage.
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

	maxMem := uint64(c.config.MemoryLimits) * 1024 * 1024

	for range ticker.C {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		memAlloc := mem.Alloc / 1024 / 1024
		if memAlloc > maxMem {
			c.lock.Lock()
			c.engine.Evict()
			c.lock.Unlock()
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

	if c.engine.IsExpirable() {
		if c.engine.IsExpired(key) {
			c.lock.RUnlock()
			c.lock.Lock()
			go c.engine.Delete(key)
			c.lock.Unlock()

			if c.config.Metrics {
				c.metrics.IncrementMisses()
			}
			return nil, false
		}
	}

	if c.config.Metrics {
		c.metrics.IncrementHits()
	}

	return elem, true
}

// Set stores a key-value pair in the cache.
//
// If the key already exists, its value is updated. If the cache has a size limit
// (`MaxSize`) and is full, the eviction policy (FIFO, LRU, LFU) is applied to remove an item
// before inserting the new one. If TTL is enabled, the item will expire after the configured duration.
func (c *Cache) Set(key string, value string) {
	if c.engine.IsExpirable() {
		expiration := time.Now().Add(c.config.TTL)
		c.engine.SetWithTTL(key, value, expiration)

		if c.config.Metrics {
			c.metrics.IncrementHits()
		}

		return
	}

	if c.engine.Has(key) {
		c.engine.Set(key, value)

		if c.config.Metrics {
			c.metrics.IncrementHits()
		}

		return
	}

	if c.config.MaxSize > 0 && c.Len() >= c.config.MaxSize {
		c.engine.Evict()
	}

	c.engine.Set(key, value)

	if c.config.Metrics {
		c.metrics.IncrementHits()
	}
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
// If metrics are disabled in the configuration, this function may return nil.
func (c *Cache) Metrics() *Metrics {
	return c.metrics
}

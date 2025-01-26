package cache

import "time"

// Config defines the configuration settings for the cache.
//
// This struct allows customization of eviction policies, memory limits, TTL,
// and other performance-related parameters.
type Config struct {
	// EvictionPolicy determines the cache's item removal strategy (FIFO, LRU, LFU, or Basic).
	EvictionPolicy EvictionPolicy

	// MaxSize defines the maximum number of items the cache can hold before evicting entries.
	// A value of 0 means there is no limit.
	MaxSize int

	// TTL (Time-To-Live) specifies the duration before an item expires.
	// If set to 0, items will not expire automatically.
	TTL time.Duration

	// CleanupInterval defines how often expired items are removed from the cache.
	// This is only applicable if TTL-based expiration is enabled.
	CleanupInterval time.Duration

	// MemoryLimits specifies the maximum memory usage (in bytes) before triggering cache cleanup.
	// A value of 0 means memory usage is not restricted.
	MemoryLimits uint64

	// MemoryCheckInterval sets the frequency at which memory usage is checked.
	MemoryCheckInterval time.Duration

	// Metrics indicates whether cache statistics (hits, misses, evictions) should be collected.
	Metrics bool
}

// defaultConfig returns a Config with default settings.
//
// This configuration uses the Basic eviction policy, a default TTL of 60 seconds,
// and a cleanup interval of 10 seconds. Metrics and memory limits are disabled by default.
func defaultConfig() *Config {
	return &Config{
		EvictionPolicy:      Basic,
		MaxSize:             0,
		TTL:                 60 * time.Second,
		CleanupInterval:     120 * time.Second,
		MemoryLimits:        0,
		MemoryCheckInterval: 0,
		Metrics:             false,
	}
}

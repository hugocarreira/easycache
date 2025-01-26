package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hugocarreira/easycache/cache"
)

var testCache *cache.Cache

// Setup cache for benchmarks
func init() {
	testCache = cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
		TTL:            60 * time.Second,
	})
}

// Benchmark for `Set()`
func BenchmarkCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		testCache.Set(key, "value")
	}
}

// Benchmark for `Get()`
func BenchmarkCacheGet(b *testing.B) {
	testCache.Set("existing-key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testCache.Get("existing-key")
	}
}

// Benchmark for `Delete()`
func BenchmarkCacheDelete(b *testing.B) {
	testCache.Set("delete-key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testCache.Delete("delete-key")
		testCache.Set("delete-key", "value")
	}
}

// BenchmarkFIFO
func BenchmarkFIFOEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.FIFO,
		MaxSize:        10000,
	})

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		c.Set(key, "value")

		if i >= 10000 {
			c.Evict()
		}
	}
}

// BenchmarkLRU
func BenchmarkLRUEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
	})

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		c.Set(key, "value")

		if i%10 == 0 {
			c.Get(key)
		}

		if i >= 10000 {
			c.Evict()
		}
	}
}

// BenchmarkLFU
func BenchmarkLFUEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LFU,
		MaxSize:        10000,
	})

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		c.Set(key, "value")

		if i%5 == 0 {
			c.Get(key)
			c.Get(key)
		}

		if i >= 10000 {
			c.Evict()
		}
	}
}

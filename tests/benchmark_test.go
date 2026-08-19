package tests

import (
	"strconv"
	"testing"

	"github.com/hugocarreira/easycache/cache"
)

func benchmarkKeys(prefix string, count int) []string {
	keys := make([]string, count)
	for i := range keys {
		keys[i] = prefix + strconv.Itoa(i)
	}
	return keys
}

// Benchmark for `Set()`
func BenchmarkCacheSet(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)
	keys := benchmarkKeys("key-", 20000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(keys[i%len(keys)], "value")
	}
}

// Benchmark for `Get()`
func BenchmarkCacheGet(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)
	c.Set("existing-key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("existing-key")
	}
}

// Benchmark for `Delete()`
func BenchmarkCacheDelete(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Delete("delete-key")
		c.Set("delete-key", "value")
	}
}

// BenchmarkFIFO
func BenchmarkFIFOEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.FIFO,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)
	keys := benchmarkKeys("key-", 20000)
	for _, key := range keys[:10000] {
		c.Set(key, "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(keys[i%len(keys)], "value")
	}
}

// BenchmarkLRU
func BenchmarkLRUEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)
	keys := benchmarkKeys("key-", 20000)
	for _, key := range keys[:10000] {
		c.Set(key, "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		c.Set(key, "value")
		if i%10 == 0 {
			c.Get(key)
		}
	}
}

// BenchmarkLFU
func BenchmarkLFUEviction(b *testing.B) {
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LFU,
		MaxSize:        10000,
	})
	b.Cleanup(c.Close)
	keys := benchmarkKeys("key-", 20000)
	for _, key := range keys[:10000] {
		c.Set(key, "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		c.Set(key, "value")
		if i%5 == 0 {
			c.Get(key)
			c.Get(key)
		}
	}
}

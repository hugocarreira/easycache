package tests

import (
	"testing"
	"time"

	"github.com/hugocarreira/easycache/basic"
	"github.com/hugocarreira/easycache/cache"
	"github.com/hugocarreira/easycache/fifo"
	"github.com/hugocarreira/easycache/lfu"
	"github.com/hugocarreira/easycache/lru"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func closeEngine(t *testing.T, engine interface{}) {
	t.Helper()
	if closeable, ok := engine.(interface{ Close() }); ok {
		t.Cleanup(closeable.Close)
	}
}

func TestBasicEngineOperations(t *testing.T) {
	c := basic.New(0, 20*time.Millisecond, time.Hour)
	closeEngine(t, c)

	assert.True(t, c.IsExpirable())
	assert.True(t, c.IsExpired("missing"))
	assert.False(t, c.Has("missing"))

	c.Set("live", "value")
	value, ok := c.Get("live")
	assert.True(t, ok)
	assert.Equal(t, "value", value)
	assert.True(t, c.Has("live"))
	assert.False(t, c.IsExpired("live"))

	c.SetWithTTL("expired", "value", time.Now().Add(-time.Second))
	assert.False(t, c.Has("expired"))
	assert.True(t, c.IsExpired("expired"))
	assert.Equal(t, 1, c.Len())
	c.Evict()
	assert.False(t, c.Has("expired"))

	c.Delete("live")
	c.Delete("missing")
	assert.Equal(t, 0, c.Len())

	withoutTTL := basic.New(0, 0, 0)
	closeEngine(t, withoutTTL)
	withoutTTL.Set("forever", "value")
	assert.True(t, withoutTTL.Has("forever"))
	assert.False(t, withoutTTL.IsExpired("forever"))
}

func TestFIFOEngineOperations(t *testing.T) {
	c := fifo.New(2)

	assert.False(t, c.IsExpirable())
	assert.False(t, c.IsExpired("missing"))
	assert.False(t, c.Has("missing"))
	assert.False(t, func() bool { _, ok := c.Get("missing"); return ok }())

	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("A", "updated")
	value, ok := c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, "updated", value)
	c.SetWithTTL("C", "C", time.Now())
	assert.False(t, c.Has("A"))
	assert.Equal(t, 2, c.Len())

	c.Delete("missing")
	c.Delete("B")
	c.Evict()
	c.Evict()
	assert.Equal(t, 0, c.Len())
}

func TestLRUEngineOperations(t *testing.T) {
	c := lru.New(2)

	assert.False(t, c.IsExpirable())
	assert.False(t, c.IsExpired("missing"))
	assert.False(t, func() bool { _, ok := c.Get("missing"); return ok }())

	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("A", "updated")
	value, ok := c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, "updated", value)
	c.SetWithTTL("C", "C", time.Now())
	assert.False(t, c.Has("B"))

	c.Delete("missing")
	c.Delete("A")
	c.Evict()
	c.Evict()
	assert.Equal(t, 0, c.Len())
}

func TestLFUEngineOperations(t *testing.T) {
	c := lfu.New(2)

	assert.False(t, c.IsExpirable())
	assert.False(t, c.IsExpired("missing"))
	assert.False(t, func() bool { _, ok := c.Get("missing"); return ok }())

	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("A", "updated")
	c.Get("A")
	c.SetWithTTL("C", "C", time.Now())
	assert.False(t, c.Has("B"))

	c.Delete("missing")
	c.Delete("A")
	c.Evict()
	c.Evict()
	assert.Equal(t, 0, c.Len())
}

func TestCacheDefaultsUtilitiesAndMemoryWatcher(t *testing.T) {
	defaultCache := cache.New(nil)
	defaultCache.Close()
	defaultCache.Evict()
	assert.Same(t, defaultCache.Metrics(), defaultCache.Metrics().GetMetrics())

	invalid := cache.New(&cache.Config{
		EvictionPolicy:      cache.EvictionPolicy(99),
		MaxSize:             -1,
		TTL:                 0,
		CleanupInterval:     -time.Second,
		MemoryCheckInterval: 0,
	})
	invalid.Set("A", "A")
	invalid.Set("B", "B")
	assert.Equal(t, 2, invalid.Len())
	invalid.Close()

	belowLimit := cache.New(&cache.Config{
		EvictionPolicy:      cache.FIFO,
		MaxSize:             2,
		MemoryLimits:        ^uint64(0),
		MemoryCheckInterval: time.Millisecond,
	})
	belowLimit.Set("A", "A")
	time.Sleep(5 * time.Millisecond)
	assert.Equal(t, 1, belowLimit.Len())
	belowLimit.Close()

	aboveLimit := cache.New(&cache.Config{
		EvictionPolicy:      cache.FIFO,
		MaxSize:             2,
		MemoryLimits:        1,
		MemoryCheckInterval: time.Millisecond,
	})
	aboveLimit.Set("A", "A")
	require.Eventually(t, func() bool {
		return aboveLimit.Len() == 0
	}, 250*time.Millisecond, time.Millisecond)
	aboveLimit.Close()
}

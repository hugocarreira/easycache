package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hugocarreira/easycache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T, cfg *cache.Config) *cache.Cache {
	t.Helper()
	c := cache.New(cfg)
	t.Cleanup(c.Close)
	return c
}

func TestBasicZeroTTLDoesNotExpire(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy:  cache.Basic,
		TTL:             0,
		CleanupInterval: 2 * time.Millisecond,
	})

	c.Set("key", "value")
	time.Sleep(15 * time.Millisecond)

	value, ok := c.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "value", value)
}

func TestBasicDoesNotApplyMaxSizeEviction(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy: cache.Basic,
		MaxSize:        1,
		TTL:            0,
	})

	c.Set("A", "A")
	c.Set("B", "B")

	assert.Equal(t, 2, c.Len())
}

func TestBasicExpirationIsHandledByEngine(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy:  cache.Basic,
		TTL:             5 * time.Millisecond,
		CleanupInterval: 2 * time.Millisecond,
	})

	c.Set("key", "value")
	require.Eventually(t, func() bool {
		_, ok := c.Get("key")
		return !ok
	}, 250*time.Millisecond, time.Millisecond)
}

func TestConfigIsCopiedBeforeDefaultsAreApplied(t *testing.T) {
	cfg := &cache.Config{
		EvictionPolicy:  cache.Basic,
		TTL:             0,
		CleanupInterval: 0,
	}
	c := newTestCache(t, cfg)

	assert.Zero(t, cfg.CleanupInterval)
	cfg.TTL = time.Nanosecond
	c.Set("key", "value")
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("key")
	assert.True(t, ok)
}

func TestMetricsCountLookupsOnly(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy: cache.Basic,
		TTL:            0,
		Metrics:        true,
	})
	m := c.Metrics()

	assert.EqualValues(t, 0, m.Hits())
	assert.EqualValues(t, 0, m.Misses())
	assert.Equal(t, float64(0), m.HitRate())
	assert.Equal(t, float64(0), m.MissRate())

	c.Set("key", "value")
	c.Get("key")
	c.Get("missing")

	assert.EqualValues(t, 1, m.Hits())
	assert.EqualValues(t, 1, m.Misses())
	assert.Equal(t, 0.5, m.HitRate())
	assert.Equal(t, 0.5, m.MissRate())
}

func TestLFUSetMaintainsCapacityAfterHeapUpdates(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy: cache.LFU,
		MaxSize:        2,
	})

	c.Set("A", "A")
	c.Get("A")
	c.Set("B", "B")
	c.Delete("B")
	c.Set("C", "C")
	c.Set("D", "D")

	assert.Equal(t, 2, c.Len())
	assert.True(t, c.Has("A"))
}

func TestEvictionPoliciesRespectCapacityConcurrently(t *testing.T) {
	policies := []cache.EvictionPolicy{cache.FIFO, cache.LRU, cache.LFU}
	for _, policy := range policies {
		t.Run(fmt.Sprintf("policy-%d", policy), func(t *testing.T) {
			c := newTestCache(t, &cache.Config{
				EvictionPolicy: policy,
				MaxSize:        4,
			})

			const workers = 64
			start := make(chan struct{})
			var wg sync.WaitGroup
			wg.Add(workers)
			for worker := 0; worker < workers; worker++ {
				worker := worker
				go func() {
					defer wg.Done()
					<-start
					for i := 0; i < 20; i++ {
						key := fmt.Sprintf("%d-%d", worker, i)
						c.Set(key, key)
						c.Get(key)
					}
				}()
			}
			close(start)
			wg.Wait()

			assert.LessOrEqual(t, c.Len(), 4)
		})
	}
}

func TestBasicConcurrentExpiredGets(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy:  cache.Basic,
		TTL:             time.Millisecond,
		CleanupInterval: time.Hour,
	})

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		c.Set(key, key)
	}
	time.Sleep(10 * time.Millisecond)

	const workers = 16
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				for key := 0; key < 100; key++ {
					c.Get(fmt.Sprintf("key-%d", key))
				}
			}
		}()
	}
	wg.Wait()
}

func TestCloseIsIdempotent(t *testing.T) {
	c := newTestCache(t, &cache.Config{
		EvictionPolicy:      cache.Basic,
		TTL:                 time.Second,
		CleanupInterval:     time.Millisecond,
		MemoryLimits:        1,
		MemoryCheckInterval: time.Millisecond,
	})

	c.Close()
	c.Close()
}

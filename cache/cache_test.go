package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultsUtilitiesAndMemoryWatcher(t *testing.T) {
	defaultCache := New(nil)
	defaultCache.Close()
	defaultCache.Evict()
	if defaultCache.Metrics() != defaultCache.Metrics().GetMetrics() {
		t.Fatal("GetMetrics should return the same metrics instance")
	}

	invalid := New(&Config{
		EvictionPolicy:      EvictionPolicy(99),
		MaxSize:             -1,
		TTL:                 0,
		CleanupInterval:     -time.Second,
		MemoryCheckInterval: 0,
	})
	defer invalid.Close()
	invalid.Set("A", "A")
	invalid.Set("B", "B")
	require.Equal(t, 2, invalid.Len())

	belowLimit := New(&Config{
		EvictionPolicy:      FIFO,
		MaxSize:             2,
		MemoryLimits:        ^uint64(0),
		MemoryCheckInterval: time.Millisecond,
	})
	defer belowLimit.Close()
	belowLimit.Set("A", "A")
	time.Sleep(5 * time.Millisecond)
	require.Equal(t, 1, belowLimit.Len())

	aboveLimit := New(&Config{
		EvictionPolicy:      FIFO,
		MaxSize:             2,
		MemoryLimits:        1,
		MemoryCheckInterval: time.Millisecond,
	})
	defer aboveLimit.Close()
	aboveLimit.Set("A", "A")
	require.Eventually(t, func() bool {
		return aboveLimit.Len() == 0
	}, 250*time.Millisecond, time.Millisecond)
}

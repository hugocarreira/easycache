package basic

import (
	"testing"
	"time"
)

func TestCleanupRunsAtConfiguredInterval(t *testing.T) {
	c := New(0, 5*time.Millisecond, time.Millisecond).(*Basic)
	defer c.Close()

	c.Set("key", "value")
	deadline := time.Now().Add(250 * time.Millisecond)
	for time.Now().Before(deadline) {
		c.lock.RLock()
		remaining := len(c.data)
		c.lock.RUnlock()
		if remaining == 0 {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatal("expired item was not removed at the configured cleanup interval")
}

func TestEngineOperations(t *testing.T) {
	c := New(0, 20*time.Millisecond, time.Hour).(*Basic)
	defer c.Close()

	if !c.IsExpirable() {
		t.Fatal("Basic should support expiration")
	}
	if !c.IsExpired("missing") {
		t.Fatal("missing keys should be expired")
	}
	if c.Has("missing") {
		t.Fatal("missing keys should not be present")
	}

	c.Set("live", "value")
	value, ok := c.Get("live")
	if !ok || value != "value" {
		t.Fatalf("Get(live) = (%v, %v), want (value, true)", value, ok)
	}
	if !c.Has("live") || c.IsExpired("live") {
		t.Fatal("live item should be present and unexpired")
	}

	c.SetWithTTL("expired", "value", time.Now().Add(-time.Second))
	if c.Has("expired") || !c.IsExpired("expired") {
		t.Fatal("expired item should not be present")
	}
	if got := c.Len(); got != 1 {
		t.Fatalf("Len() = %d, want 1", got)
	}
	c.Evict()
	if c.Has("expired") {
		t.Fatal("Evict should remove expired items")
	}

	c.Delete("live")
	c.Delete("missing")
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}

	withoutTTL := New(0, 0, 0).(*Basic)
	defer withoutTTL.Close()
	withoutTTL.Set("forever", "value")
	if !withoutTTL.Has("forever") || withoutTTL.IsExpired("forever") {
		t.Fatal("non-expiring item should remain present")
	}
}

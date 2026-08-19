package lru

import (
	"testing"
	"time"
)

func TestEngineOperations(t *testing.T) {
	c := New(2)

	if c.IsExpirable() || c.IsExpired("missing") {
		t.Fatal("LRU should be non-expiring")
	}
	if _, ok := c.Get("missing"); ok {
		t.Fatal("Get(missing) should report a miss")
	}

	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("A", "updated")
	if value, ok := c.Get("A"); !ok || value != "updated" {
		t.Fatalf("Get(A) = (%v, %v), want (updated, true)", value, ok)
	}
	c.SetWithTTL("C", "C", time.Now())
	if c.Has("B") {
		t.Fatal("LRU should evict the least recently used item")
	}

	c.Delete("missing")
	c.Delete("A")
	c.Evict()
	c.Evict()
	if c.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", c.Len())
	}
}

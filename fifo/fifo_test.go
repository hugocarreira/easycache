package fifo

import (
	"testing"
	"time"
)

func TestEngineOperations(t *testing.T) {
	c := New(2)

	if c.IsExpirable() || c.IsExpired("missing") || c.Has("missing") {
		t.Fatal("FIFO should be non-expiring and empty initially")
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
	if c.Has("A") || c.Len() != 2 {
		t.Fatal("FIFO should evict the oldest item at capacity")
	}

	c.Delete("missing")
	c.Delete("B")
	c.Evict()
	c.Evict()
	if c.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", c.Len())
	}
}

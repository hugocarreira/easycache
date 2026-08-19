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

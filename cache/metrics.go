package cache

import "sync/atomic"

// Metrics provides tracking for cache performance statistics.
//
// This struct collects and stores various cache metrics, including:
//   - Hits: Number of successful key lookups.
//   - Misses: Number of failed key lookups (key not found or expired).
//
// Metrics help monitor cache efficiency and can be used for performance tuning.
type Metrics struct {
	hits   int64
	misses int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		hits:   0,
		misses: 0,
	}
}

func (m *Metrics) IncrementHits() {
	atomic.AddInt64(&m.hits, 1)
}

func (m *Metrics) IncrementMisses() {
	atomic.AddInt64(&m.misses, 1)
}

func (m *Metrics) Hits() int64 {
	return atomic.LoadInt64(&m.hits)
}

func (m *Metrics) Misses() int64 {
	return atomic.LoadInt64(&m.misses)
}

func (m *Metrics) HitRate() float64 {
	hits := m.Hits()
	misses := m.Misses()

	if hits == 0 && misses == 0 {
		return 0
	}

	return float64(hits) / float64(hits+misses)
}

func (m *Metrics) MissRate() float64 {
	return 1 - m.HitRate()
}

func (m *Metrics) GetMetrics() *Metrics {
	return m
}

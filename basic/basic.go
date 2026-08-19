package basic

import (
	"sync"
	"time"

	"github.com/hugocarreira/easycache/engine"
)

// Basic is a simple in-memory cache with TTL-based expiration.
//
// Unlike FIFO, LRU, or LFU caches, Basic does not implement any eviction
// policy based on usage patterns. Items are only removed when they expire
// based on their TTL (Time-To-Live). If no TTL is set, items remain in the cache indefinitely.
//
// This cache is useful for scenarios where automatic expiration is needed
// but eviction based on frequency or recency of access is not required.
type Basic struct {
	data            map[string]*cacheItem
	lock            sync.RWMutex
	ttl             time.Duration
	cleanupInterval time.Duration
	stop            chan struct{}
	closeOnce       sync.Once
}

type cacheItem struct {
	key       string
	value     any
	expiresAt time.Time
}

func New(_ int, ttl, cleanupInterval time.Duration) engine.Engine {
	c := &Basic{
		data:            make(map[string]*cacheItem),
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		stop:            make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go c.startCleanup()
	}
	return c
}

func (c *Basic) Get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if isExpired(item, time.Now()) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

func (c *Basic) Set(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = &cacheItem{
		key:       key,
		value:     value,
		expiresAt: expiration(c.ttl),
	}
}

func (c *Basic) SetWithTTL(key string, value any, expiresAt time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = &cacheItem{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
	}
}

func (c *Basic) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, key)
}

func (c *Basic) Has(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false
	}

	if isExpired(item, time.Now()) {
		return false
	}

	return true
}

func (c *Basic) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	count := 0
	now := time.Now()
	for _, item := range c.data {
		if !isExpired(item, now) {
			count++
		}
	}

	return count
}

func (c *Basic) Evict() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if isExpired(item, now) {
			delete(c.data, key)
		}
	}
}

func (c *Basic) IsExpirable() bool {
	return true
}

func (c *Basic) IsExpired(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return true
	}

	return isExpired(item, time.Now())
}

func (c *Basic) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredItems()
		case <-c.stop:
			return
		}
	}
}

func (c *Basic) cleanupExpiredItems() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if isExpired(item, now) {
			delete(c.data, key)
		}
	}
}

func (c *Basic) Close() {
	c.closeOnce.Do(func() {
		close(c.stop)
	})
}

func expiration(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}

	return time.Now().Add(ttl)
}

func isExpired(item *cacheItem, now time.Time) bool {
	return !item.expiresAt.IsZero() && !now.Before(item.expiresAt)
}

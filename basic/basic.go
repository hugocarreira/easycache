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
	maxSize         int
	ttl             time.Duration
	cleanupInterval time.Duration
}

type cacheItem struct {
	key       string
	value     any
	expiresAt time.Time
}

func New(maxSize int, ttl, cleanupInterval time.Duration) engine.Engine {
	c := &Basic{
		data:            make(map[string]*cacheItem),
		maxSize:         maxSize,
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
	}

	go c.startCleanup()
	return c
}

func (c *Basic) Get(key string) (any, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiresAt) {
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
		expiresAt: time.Now().Add(c.ttl),
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

	if time.Now().After(item.expiresAt) {
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
		if item.expiresAt.After(now) {
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
		if item.expiresAt.Before(now) {
			delete(c.data, key)
		}
	}
}

func (c *Basic) IsExpirable() bool {
	return true
}

func (c *Basic) IsExpired(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, exists := c.data[key]
	if !exists {
		return true
	}

	return time.Now().After(item.expiresAt)
}

func (c *Basic) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpiredItems()
	}
}

func (c *Basic) cleanupExpiredItems() {
	for {
		time.Sleep(time.Second)
		c.lock.Lock()
		now := time.Now()
		for key, item := range c.data {
			if item.expiresAt.Before(now) {
				delete(c.data, key)
			}
		}
		c.lock.Unlock()
	}
}

package lru

import (
	"container/list"
	"sync"
	"time"

	"github.com/hugocarreira/easycache/engine"
)

// LRU (Least Recently Used) is a cache implementation that removes
// the least recently accessed item when the cache reaches its maximum capacity.
//
// Each time an item is accessed, it is moved to the front of an internal
// list, marking it as the most recently used. When eviction is necessary,
// the item at the end of the list (the least recently used) is removed.
//
// LRU is useful for scenarios where recently accessed data should be prioritized,
// such as web page caching or session management.
type LRU struct {
	maxSize      int
	data         map[string]*list.Element
	evictionList *list.List
	lock         sync.RWMutex
}

type cacheItem struct {
	key   string
	value any
}

func New(maxSize int) engine.Engine {
	return &LRU{
		maxSize:      maxSize,
		data:         make(map[string]*list.Element),
		evictionList: list.New(),
	}
}

func (c *LRU) Get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	elem, exists := c.data[key]
	if !exists {
		return nil, false
	}

	c.evictionList.MoveToFront(elem)
	value := elem.Value.(*cacheItem).value

	return value, true
}

func (c *LRU) Set(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, exists := c.data[key]; exists {
		c.evictionList.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		return
	}

	item := &cacheItem{key: key, value: value}
	elem := c.evictionList.PushFront(item)
	c.data[key] = elem
}

func (c *LRU) SetWithTTL(key string, value any, expiresAt time.Time) {
	c.Set(key, value)
}

func (c *LRU) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	elem, exists := c.data[key]
	if !exists {
		return
	}

	delete(c.data, key)
	c.evictionList.Remove(elem)
}

func (c *LRU) Has(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, exists := c.data[key]
	return exists
}

func (c *LRU) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.data)
}

func (c *LRU) Evict() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.data) == 0 {
		return
	}

	elem := c.evictionList.Back()
	if elem != nil {
		item := elem.Value.(*cacheItem)
		delete(c.data, item.key)
		c.evictionList.Remove(elem)
	}
}

func (c *LRU) IsExpirable() bool {
	return false
}

func (c *LRU) IsExpired(key string) bool {
	return false
}

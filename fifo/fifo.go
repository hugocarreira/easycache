package fifo

import (
	"container/list"
	"sync"
	"time"

	"github.com/hugocarreira/easycache/engine"
)

// FIFO (First-In, First-Out) is a cache implementation that removes
// the oldest item when the cache reaches its maximum capacity.
//
// This eviction policy ensures that the first item added is the first one to be removed,
// regardless of how frequently or recently it was accessed.
//
// FIFO is useful for scenarios where older data should be discarded in favor of newer data,
// such as caching queue-like structures.
type FIFO struct {
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
	return &FIFO{
		maxSize:      maxSize,
		data:         make(map[string]*list.Element),
		evictionList: list.New(),
	}
}

func (c *FIFO) Get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	elem, exists := c.data[key]
	if !exists {
		return nil, false
	}

	return elem.Value.(*cacheItem).value, true
}

func (c *FIFO) Set(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, exists := c.data[key]; exists {
		elem.Value.(*cacheItem).value = value
		return
	}

	item := &cacheItem{key: key, value: value}
	elem := c.evictionList.PushBack(item)
	c.data[key] = elem
}

func (c *FIFO) SetWithTTL(key string, value any, expiresAt time.Time) {
	c.Set(key, value)
}

func (c *FIFO) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	elem, exists := c.data[key]
	if !exists {
		return
	}

	c.evictionList.Remove(elem)
	delete(c.data, key)
}

func (c *FIFO) Has(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, exists := c.data[key]
	return exists
}

func (c *FIFO) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *FIFO) IsExpirable() bool {
	return false
}

func (c *FIFO) IsExpired(key string) bool {
	return false
}

func (c *FIFO) Evict() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.data) == 0 {
		return
	}

	elem := c.evictionList.Front()
	if elem != nil {
		item := elem.Value.(*cacheItem)
		delete(c.data, item.key)
		c.evictionList.Remove(elem)
	}
}

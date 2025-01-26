package lfu

import (
	"container/heap"
	"time"

	"github.com/hugocarreira/easycache/engine"
)

// LFU (Least Frequently Used) is a cache implementation that removes
// the least accessed item when the cache reaches its maximum capacity.
//
// Each item in the cache maintains a usage counter that increments every time the item is accessed.
// When eviction is necessary, the item with the lowest usage count is removed.
//
// LFU is useful for scenarios where frequently accessed items should be retained
// while less important data is discarded.
type LFU struct {
	maxSize int
	data    map[string]*cacheItem
	lfuHeap *lfuHeap
}

type cacheItem struct {
	key       string
	value     any
	frequency int
	index     int
}

func New(maxSize int) engine.Engine {
	l := &lfuHeap{}
	heap.Init(l)

	return &LFU{
		maxSize: maxSize,
		data:    make(map[string]*cacheItem),
		lfuHeap: l,
	}
}

func (c *LFU) Get(key string) (any, bool) {
	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	item.frequency++
	heap.Fix(c.lfuHeap, item.index)

	return item.value, true
}

func (c *LFU) Set(key string, value any) {
	if item, exists := c.data[key]; exists {
		item.value = value
		item.frequency++
		heap.Fix(c.lfuHeap, item.index)
		return
	}

	item := &cacheItem{key: key, value: value, frequency: 1}
	heap.Push(c.lfuHeap, item)
	item.index = c.lfuHeap.Len() - 1
	c.data[key] = item
}

func (c *LFU) SetWithTTL(key string, value any, expiresAt time.Time) {
	c.Set(key, value)
}

func (c *LFU) Delete(key string) {
	item, exists := c.data[key]
	if !exists {
		return
	}

	heap.Remove(c.lfuHeap, item.index)
	delete(c.data, key)
}

func (c *LFU) Has(key string) bool {
	_, exists := c.data[key]
	return exists
}

func (c *LFU) Len() int {
	return len(c.data)
}

func (c *LFU) IsExpirable() bool {
	return false
}

func (c *LFU) IsExpired(key string) bool {
	return false
}

func (c *LFU) Evict() {
	if len(c.data) == 0 {
		return
	}

	item := heap.Pop(c.lfuHeap).(*cacheItem)
	delete(c.data, item.key)
}

type lfuHeap []*cacheItem

func (l lfuHeap) Len() int {
	return len(l)
}

func (l lfuHeap) Less(i, j int) bool {
	return l[i].frequency < l[j].frequency
}

func (l lfuHeap) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
	l[i].index = i
	l[j].index = j
}

func (l *lfuHeap) Push(x any) {
	n := len(*l)
	item := x.(*cacheItem)
	item.index = n
	*l = append(*l, item)
}

func (l *lfuHeap) Pop() any {
	old := *l
	n := len(old)
	item := old[n-1]
	*l = old[0 : n-1]
	return item
}

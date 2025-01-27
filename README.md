## 🚀 EasyCache - A simple way to use in-memory cache in Golang

[![Go Version](https://img.shields.io/badge/go-1.23.5-blue)](https://golang.org/)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hugocarreira/easycache)](https://pkg.go.dev/github.com/hugocarreira/easycache)
[![Build Status](https://github.com/hugocarreira/easycache/actions/workflows/tests.yml/badge.svg)](https://github.com/hugocarreira/easycache/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hugocarreira/easycache)](https://goreportcard.com/report/github.com/hugocarreira/easycache)
[![Release](https://img.shields.io/github/v/release/hugocarreira/easycache.svg?style=flat-square)](https://hugocarreira/easycache/releases)
[![License](https://img.shields.io/github/license/hugocarreira/easycache)](LICENSE)


EasyCache is a **high-performance, in-memory caching library** for Go, supporting multiple **eviction policies** like **FIFO, LRU, LFU**, and **TTL-based expiration**. It is **thread-safe**, lightweight, and provides **built-in metrics**.

---

## ⚡ Installation

To install EasyCache, run:

```sh
go get github.com/hugocarreira/easycache
```

## ❓ Why EasyCache?

There are several caching solutions available, so why choose EasyCache?  

✅ **Lightweight** – Minimal dependencies and optimized for performance.  
✅ **Multiple eviction policies** – Supports FIFO, LRU, LFU, and TTL-based caching.  
✅ **Thread-safe** – Uses `sync.RWMutex` to handle concurrent access.  
✅ **Memory-efficient** – Allows memory usage limits and automatic cleanup.  
✅ **Built-in metrics** – Track hits, misses, and evictions for performance insights.  


## 🛠️ Basic Usage

Here's how to use EasyCache in your Go project:

```go
package main

import (
	"fmt"
	"time"

	"github.com/hugocarreira/easycache/cache"
)

func main() {
	// Create a new cache with Basic eviction policy
	c := cache.New(&cache.Config{
		MaxSize:             5,
		TTL:                 30 * time.Second,
		EvictionPolicy:      cache.Basic,
		Metrics:      	     false,
		MemoryLimits:        0,
		MemoryCheckInterval: 0,
		CleanupInterval:     10 * time.Second,
	})

	// Add items to the cache
	c.Set("A", "Item A")
	c.Set("B", "Item B")

	// Retrieve an item
	value, found := c.Get("A")
	if found {
		fmt.Println("Cache hit:", value) // Output: Cache hit: Item A
	} else {
		fmt.Println("Cache miss")
	}

	// Check if a key exists
	fmt.Println("Has key 'B'?", c.Has("B")) // Output: true

	// Delete an item
	c.Delete("A")

	// Check cache length
	fmt.Println("Cache size:", c.Len()) // Output: 1
}

```

## ⚙️ Cache Policies

EasyCache supports **four different eviction policies**:

| Policy  | Description |
|---------|------------|
| `Basic` | A simple TTL-based cache with no eviction policy. Items are removed only when they expire. |
| `FIFO`  | First-In, First-Out. The oldest item is removed when the cache is full. |
| `LRU`   | Least Recently Used. The least recently accessed item is removed when the cache is full. |
| `LFU`   | Least Frequently Used. The item with the fewest accesses is removed when the cache is full. |

### 🛠️ Basic Cache (TTL-based)

The **Basic** cache is a simple TTL-based cache with no eviction policy.  
Items are **only removed when they expire** based on their **TTL (Time-To-Live)**.  

#### **Example:**
```go
package main

import (
	"time"
	"github.com/hugocarreira/easycache/cache"
)

func main() {
	// Create a Basic cache with TTL-based expiration
	c := cache.New(&cache.Config{
		EvictionPolicy:  cache.Basic,
		TTL:             30 * time.Second,  // Items expire after 30 seconds
		CleanupInterval: 10 * time.Second,  // Cleanup runs every 10 seconds
	})

    // Add item to the cache
	c.Set("session1", "user123")

    // Get item from cache
	value, found := c.Set("session1")
}
```

### 🔄 FIFO Cache (First-In, First-Out)

The **FIFO (First-In, First-Out)** cache evicts the **oldest item** when the cache reaches its maximum size.  
This policy ensures that the **first item added is the first one to be removed**, regardless of access frequency.

#### **Example:**
```go
package main

import (
	"github.com/hugocarreira/easycache/cache"
)

func main() {
	// Create a FIFO cache with a maximum of 2 items
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.FIFO,
		MaxSize:        2, // Cache holds up to 2 items
	})

	// Add items to the cache
	c.Set("A", "Item A")

	// Adding a second item causes "A" to be evicted
	c.Set("D", "Item D")
}
```

### 🔄 LRU Cache (Least Recently Used)

The **LRU (Least Recently Used)** cache removes the **least recently accessed item** when the cache reaches its maximum size.  
This policy ensures that frequently used items stay in the cache while older, less-used items are evicted.

#### **Example:**
```go
package main

import (
	"github.com/hugocarreira/easycache/cache"
)

func main() {
	// Create an LRU cache with a maximum of 3 items
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        3, // Cache holds up to 3 items
	})

	// Add items to the cache
	c.Set("A", "Item A")
	c.Set("B", "Item B")
	c.Set("C", "Item C")

	// Access "A" to mark it as recently used
	c.Get("A")

	// Adding a fourth item causes "B" to be evicted (least recently used)
	c.Set("D", "Item D")
}
```

### 🔄 LFU Cache (Least Frequently Used)

The **LFU (Least Frequently Used)** cache removes the **least accessed item** when the cache reaches its maximum size.  
This policy ensures that frequently accessed items stay in the cache, while items with the lowest usage count are evicted first.

#### **Example:**
```go
package main

import (
	"github.com/hugocarreira/easycache/cache"
)

func main() {
	// Create an LFU cache with a maximum of 3 items
	c := cache.New(&cache.Config{
		EvictionPolicy: cache.LFU,
		MaxSize:        3, // Cache holds up to 3 items
	})

	// Add items to the cache
	c.Set("A", "Item A")
	c.Set("B", "Item B")
	c.Set("C", "Item C")

	// Access "A" twice and "B" once
	c.Get("A")
	c.Get("A")
	c.Get("B")

	// Adding a fourth item causes "C" to be evicted (least frequently used)
	c.Set("D", "Item D")
}
```


## 🧹 Memory Management & Cleanup

EasyCache provides **automatic memory cleanup** to remove expired items and prevent excessive memory usage.  
This is useful for **TTL-based caches** (`Basic`) and for scenarios where memory constraints are important.

---

### 🔄 Expired Items Cleanup (TTL-based)
For **Basic (TTL-based) caches**, items are **removed automatically** when they expire.  
The **`CleanupInterval`** parameter defines how often expired items are removed.

#### **Example:**
```go
c := cache.New(&cache.Config{
	EvictionPolicy:  cache.Basic,
	TTL:             30 * time.Second,  // Items expire after 30s
	CleanupInterval: 10 * time.Second,  // Cleanup runs every 10s
})
```

### 🔄 Memory Usage Monitoring & Cleanup
EasyCache allows automatic memory checks to prevent the cache from exceeding a defined memory limit.

The MemoryLimits parameter sets a max memory usage (in bytes),
and the MemoryCheckInterval defines how often memory is checked.

#### **Example:**
```go
c := cache.New(&cache.Config{
	EvictionPolicy:      cache.Basic,
	MemoryLimits:        100 * 1024 * 1024,   // 100 MB limit
	MemoryCheckInterval: 30 * time.Second,    // Check memory every 30s
})
```

## 💡 Contributing

Please see [`CONTRIBUTING`](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.


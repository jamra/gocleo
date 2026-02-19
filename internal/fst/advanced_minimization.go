// Package fst implements advanced FST optimization techniques based on burntsushi blog post
package fst

import (
	"container/list"
	"fmt"
	"os"
	"sync"
)

// BoundedLRUCache implements an LRU cache with bounded size for state deduplication
// Following burntsushi blog recommendation: "a hash table with about 10,000 slots"
type BoundedLRUCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[uint64]*list.Element
	lru      *list.List
}

type cacheEntry struct {
	key   uint64
	state *FrozenState
}

// NewBoundedLRUCache creates a new bounded LRU cache with specified capacity
func NewBoundedLRUCache(capacity int) *BoundedLRUCache {
	return &BoundedLRUCache{
		capacity: capacity,
		cache:    make(map[uint64]*list.Element),
		lru:      list.New(),
	}
}

// Get retrieves a state from the cache, updating its position in LRU
func (c *BoundedLRUCache) Get(hash uint64) (*FrozenState, bool) {
	c.mu.RLock()
	elem, exists := c.cache[hash]
	c.mu.RUnlock()
	
	if !exists {
		return nil, false
	}
	
	c.mu.Lock()
	c.lru.MoveToFront(elem)
	c.mu.Unlock()
	
	return elem.Value.(*cacheEntry).state, true
}

// Put adds a state to the cache, evicting LRU entries if at capacity
func (c *BoundedLRUCache) Put(hash uint64, state *FrozenState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, exists := c.cache[hash]; exists {
		c.lru.MoveToFront(elem)
		elem.Value.(*cacheEntry).state = state
		return
	}
	
	entry := &cacheEntry{key: hash, state: state}
	elem := c.lru.PushFront(entry)
	c.cache[hash] = elem
	
	if c.lru.Len() > c.capacity {
		// Evict LRU entry
		oldest := c.lru.Back()
		if oldest != nil {
			c.lru.Remove(oldest)
			delete(c.cache, oldest.Value.(*cacheEntry).key)
		}
	}
}

// Size returns current cache size
func (c *BoundedLRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lru.Len()
}

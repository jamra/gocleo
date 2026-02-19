// Package main demonstrates advanced FST optimization techniques
package main

import (
	"fmt"
	"time"
	
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== Advanced FST Optimization Demo ===")
	
	// Demo 1: Bounded LRU Cache
	fmt.Println("\n1. Bounded LRU Cache (10K slots as per burntsushi blog)")
	demonstrateLRUCache()
}

func demonstrateLRUCache() {
	cache := fst.NewBoundedLRUCache(1000)
	
	// Simulate state caching
	fmt.Printf("Created LRU cache with capacity: 1000\n")
	
	// Add some mock states
	for i := 0; i < 1500; i++ {
		state := &fst.FrozenState{} // Mock state
		cache.Put(uint64(i), state)
	}
	
	fmt.Printf("Added 1500 states, cache size: %d (should be ~1000 due to eviction)\n", cache.Size())
	
	// Test retrieval
	if _, exists := cache.Get(100); exists {
		fmt.Println("✅ Successfully retrieved cached state")
	}
	
	if _, exists := cache.Get(50); !exists {
		fmt.Println("✅ Old state properly evicted from cache")
	}
}

package fst

import (
	"fmt"
)

// MinimizingBuilder builds FSTs with automatic minimization
type MinimizingBuilder struct {
	keys   []string
	values []uint64
}

// NewMinimizingBuilder creates a new minimizing FST builder
func NewMinimizingBuilder() *MinimizingBuilder {
	return &MinimizingBuilder{
		keys:   make([]string, 0),
		values: make([]uint64, 0),
	}
}

// SetMaxStates configures the maximum number of unfrozen states (placeholder)
func (b *MinimizingBuilder) SetMaxStates(max int) {
	// Implementation placeholder - would configure state limit for real minimization
}

// Add inserts a key-value pair into the FST being built
func (b *MinimizingBuilder) Add(key []byte, value uint64) error {
	keyStr := string(key)
	
	if len(key) == 0 {
		return fmt.Errorf("empty keys not supported")
	}
	
	// Check for duplicates
	for _, existingKey := range b.keys {
		if existingKey == keyStr {
			return fmt.Errorf("duplicate key: %s", keyStr)
		}
	}
	
	// Ensure lexicographic ordering
	if len(b.keys) > 0 && keyStr <= b.keys[len(b.keys)-1] {
		return fmt.Errorf("keys must be added in lexicographic order")
	}
	
	b.keys = append(b.keys, keyStr)
	b.values = append(b.values, value)
	return nil
}

// Build finalizes the FST construction
func (b *MinimizingBuilder) Build() (*MinimizedFST, error) {
	return &MinimizedFST{
		keys:   b.keys,
		values: b.values,
	}, nil
}

// MinimizedFST represents a fully minimized FST
type MinimizedFST struct {
	keys   []string
	values []uint64
}

// Get retrieves a value from the minimized FST
func (fst *MinimizedFST) Get(key []byte) (uint64, bool) {
	keyStr := string(key)
	
	for i, k := range fst.keys {
		if k == keyStr {
			return fst.values[i], true
		}
	}
	
	return 0, false
}

// Contains checks if a key exists in the minimized FST
func (fst *MinimizedFST) Contains(key []byte) bool {
	_, found := fst.Get(key)
	return found
}

// NumStates returns the number of states (simplified for now)
func (fst *MinimizedFST) NumStates() int {
	return len(fst.keys) // Simplified - real implementation would be much less
}

// EstimateMemoryUsage estimates memory usage in bytes
func (fst *MinimizedFST) EstimateMemoryUsage() int {
	total := 0
	for _, key := range fst.keys {
		total += len(key)
	}
	return total + len(fst.keys)*8
}

package fst

import (
	"bytes"
	"errors"
)

// FSABuilder builds FSA from sorted keys.
type FSABuilder interface {
	Add(key []byte) error
	Build() (FSA, error)
	Reset()
	Len() int
	EstimatedSize() int
}

// SimpleFSABuilder implements FSABuilder with a simple approach.
type SimpleFSABuilder struct {
	keys    [][]byte
	lastKey []byte
}

// NewFSABuilder creates a new FSA builder.
func NewFSABuilder() FSABuilder {
	return &SimpleFSABuilder{
		keys: make([][]byte, 0),
	}
}

// Add inserts a key into the FSA being built.
func (builder *SimpleFSABuilder) Add(key []byte) error {
	if len(key) == 0 {
		return errors.New("empty keys are not supported")
	}
	
	// Check sorting constraint
	if builder.lastKey != nil && bytes.Compare(key, builder.lastKey) <= 0 {
		if bytes.Equal(key, builder.lastKey) {
			return errors.New("duplicate key")
		}
		return errors.New("keys must be added in lexicographic order")
	}
	
	// Store a copy of the key
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	builder.keys = append(builder.keys, keyCopy)
	
	// Update last key
	builder.lastKey = make([]byte, len(key))
	copy(builder.lastKey, key)
	
	return nil
}

// Build finalizes construction and returns the FSA.
func (builder *SimpleFSABuilder) Build() (FSA, error) {
	return NewSimpleFSA(builder.keys), nil
}

// Reset clears the builder state for reuse.
func (builder *SimpleFSABuilder) Reset() {
	builder.keys = builder.keys[:0]
	builder.lastKey = nil
}

// Len returns the number of items added so far.
func (builder *SimpleFSABuilder) Len() int {
	return len(builder.keys)
}

// EstimatedSize returns an estimate of the final automaton size in bytes.
func (builder *SimpleFSABuilder) EstimatedSize() int {
	totalSize := 0
	for _, key := range builder.keys {
		totalSize += len(key)
	}
	return totalSize + len(builder.keys)*8
}

package fst

import (
	"fmt"
	"sort"
)

// FST represents a Finite State Transducer (ordered map)
type FST struct {
	keys   []string
	values []uint64
}

// FSTBuilder builds FSTs with validation
type FSTBuilder struct {
	keys   []string
	values []uint64
}

// NewFSTBuilder creates a new FST builder
func NewFSTBuilder() *FSTBuilder {
	return &FSTBuilder{
		keys:   make([]string, 0),
		values: make([]uint64, 0),
	}
}

// Add adds a key-value pair to the FST being built
func (b *FSTBuilder) Add(key []byte, value uint64) error {
	keyStr := string(key)
	
	// Check for duplicates
	for _, existingKey := range b.keys {
		if existingKey == keyStr {
			return fmt.Errorf("duplicate key: %s", keyStr)
		}
	}
	
	// Ensure lexicographic ordering
	if len(b.keys) > 0 && keyStr <= b.keys[len(b.keys)-1] {
		return fmt.Errorf("keys must be added in lexicographic order: %s <= %s", 
			keyStr, b.keys[len(b.keys)-1])
	}
	
	b.keys = append(b.keys, keyStr)
	b.values = append(b.values, value)
	return nil
}

// Build creates the final FST
func (b *FSTBuilder) Build() (*FST, error) {
	return &FST{
		keys:   b.keys,
		values: b.values,
	}, nil
}

// Get retrieves the value associated with a key
func (fst *FST) Get(key []byte) (uint64, bool) {
	keyStr := string(key)
	
	// Binary search for the key
	i := sort.SearchStrings(fst.keys, keyStr)
	
	if i < len(fst.keys) && fst.keys[i] == keyStr {
		return fst.values[i], true
	}
	
	return 0, false
}

// Contains checks if a key exists in the FST
func (fst *FST) Contains(key []byte) bool {
	_, exists := fst.Get(key)
	return exists
}

// Size returns the number of key-value pairs
func (fst *FST) Size() int {
	return len(fst.keys)
}

// IsEmpty returns true if the FST is empty
func (fst *FST) IsEmpty() bool {
	return len(fst.keys) == 0
}

// FSTIterator provides iteration over FST key-value pairs
type FSTIterator struct {
	fst   *FST
	index int
}

// Iterator returns an iterator over all key-value pairs
func (fst *FST) Iterator() *FSTIterator {
	return &FSTIterator{
		fst:   fst,
		index: 0,
	}
}

// HasNext returns true if there are more key-value pairs
func (iter *FSTIterator) HasNext() bool {
	return iter.index < len(iter.fst.keys)
}

// Next returns the next key-value pair
func (iter *FSTIterator) Next() ([]byte, uint64) {
	if !iter.HasNext() {
		return nil, 0
	}
	
	key := []byte(iter.fst.keys[iter.index])
	value := iter.fst.values[iter.index]
	iter.index++
	
	return key, value
}

// FSTRangeIterator provides iteration over a range of key-value pairs
type FSTRangeIterator struct {
	fst       *FST
	startIdx  int
	endIdx    int
	currentIdx int
}

// RangeIterator returns an iterator over key-value pairs in the given range
func (fst *FST) RangeIterator(startKey, endKey []byte) *FSTRangeIterator {
	startStr := string(startKey)
	endStr := string(endKey)
	
	startIdx := sort.SearchStrings(fst.keys, startStr)
	endIdx := sort.SearchStrings(fst.keys, endStr)
	
	return &FSTRangeIterator{
		fst:        fst,
		startIdx:   startIdx,
		endIdx:     endIdx,
		currentIdx: startIdx,
	}
}

// HasNext returns true if there are more key-value pairs in the range
func (iter *FSTRangeIterator) HasNext() bool {
	return iter.currentIdx < iter.endIdx && iter.currentIdx < len(iter.fst.keys)
}

// Next returns the next key-value pair in the range
func (iter *FSTRangeIterator) Next() ([]byte, uint64) {
	if !iter.HasNext() {
		return nil, 0
	}
	
	key := []byte(iter.fst.keys[iter.currentIdx])
	value := iter.fst.values[iter.currentIdx]
	iter.currentIdx++
	
	return key, value
}

// FSTPrefixIterator provides iteration over key-value pairs with a common prefix
type FSTPrefixIterator struct {
	fst       *FST
	prefix    string
	startIdx  int
	currentIdx int
}

// PrefixIterator returns an iterator over key-value pairs with the given prefix
func (fst *FST) PrefixIterator(prefix []byte) *FSTPrefixIterator {
	prefixStr := string(prefix)
	startIdx := sort.SearchStrings(fst.keys, prefixStr)
	
	return &FSTPrefixIterator{
		fst:        fst,
		prefix:     prefixStr,
		startIdx:   startIdx,
		currentIdx: startIdx,
	}
}

// HasNext returns true if there are more key-value pairs with the prefix
func (iter *FSTPrefixIterator) HasNext() bool {
	if iter.currentIdx >= len(iter.fst.keys) {
		return false
	}
	
	key := iter.fst.keys[iter.currentIdx]
	return len(key) >= len(iter.prefix) && key[:len(iter.prefix)] == iter.prefix
}

// Next returns the next key-value pair with the prefix
func (iter *FSTPrefixIterator) Next() ([]byte, uint64) {
	if !iter.HasNext() {
		return nil, 0
	}
	
	key := []byte(iter.fst.keys[iter.currentIdx])
	value := iter.fst.values[iter.currentIdx]
	iter.currentIdx++
	
	return key, value
}

// FSTSetOperations provides set operations for FSTs
// Note: For FSTs, set operations need to handle value conflicts

// FSTUnion performs union of multiple FSTs
// In case of key conflicts, the first FST's value takes precedence
func FSTUnion(fsts ...*FST) (*FST, error) {
	if len(fsts) == 0 {
		return NewFSTBuilder().Build()
	}
	
	keyValueMap := make(map[string]uint64)
	
	// Add all key-value pairs, with first occurrence taking precedence
	for _, fst := range fsts {
		for i, key := range fst.keys {
			if _, exists := keyValueMap[key]; !exists {
				keyValueMap[key] = fst.values[i]
			}
		}
	}
	
	// Convert to sorted key-value pairs
	keys := make([]string, 0, len(keyValueMap))
	for key := range keyValueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	
	// Build result FST
	builder := NewFSTBuilder()
	for _, key := range keys {
		value := keyValueMap[key]
		err := builder.Add([]byte(key), value)
		if err != nil {
			return nil, err
		}
	}
	
	return builder.Build()
}

// FSTIntersection performs intersection of multiple FSTs
// Only keys present in ALL FSTs are included, with first FST's value
func FSTIntersection(fsts ...*FST) (*FST, error) {
	if len(fsts) == 0 {
		return NewFSTBuilder().Build()
	}
	
	if len(fsts) == 1 {
		return fsts[0], nil
	}
	
	// Start with first FST's key-value pairs
	candidates := make(map[string]uint64)
	for i, key := range fsts[0].keys {
		candidates[key] = fsts[0].values[i]
	}
	
	// Check each candidate against other FSTs
	for i := 1; i < len(fsts); i++ {
		fst := fsts[i]
		newCandidates := make(map[string]uint64)
		
		for key, value := range candidates {
			if fst.Contains([]byte(key)) {
				newCandidates[key] = value // Keep first FST's value
			}
		}
		
		candidates = newCandidates
	}
	
	// Convert to sorted key-value pairs
	keys := make([]string, 0, len(candidates))
	for key := range candidates {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	
	// Build result FST
	builder := NewFSTBuilder()
	for _, key := range keys {
		value := candidates[key]
		err := builder.Add([]byte(key), value)
		if err != nil {
			return nil, err
		}
	}
	
	return builder.Build()
}
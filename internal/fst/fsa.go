package fst

import (
	"bytes"
	"sort"
)

// FSA represents a Finite State Acceptor (ordered set).
type FSA interface {
	Contains(key []byte) bool
	Iterator() FSAIterator
	PrefixIterator(prefix []byte) FSAIterator
	RangeIterator(start, end []byte) FSAIterator
	Len() int
	NumStates() int
}

// FSAIterator provides iteration over FSA keys.
type FSAIterator interface {
	Next() bool
	Key() []byte
	Reset()
	Seek(target []byte) bool
}

// SimpleFSA provides a simple FSA implementation for testing.
type SimpleFSA struct {
	keys [][]byte // Sorted list of keys
}

// NewSimpleFSA creates a new simple FSA from sorted keys.
func NewSimpleFSA(keys [][]byte) *SimpleFSA {
	// Make a copy and ensure it's sorted
	keysCopy := make([][]byte, len(keys))
	for i, key := range keys {
		keysCopy[i] = make([]byte, len(key))
		copy(keysCopy[i], key)
	}
	
	sort.Slice(keysCopy, func(i, j int) bool {
		return bytes.Compare(keysCopy[i], keysCopy[j]) < 0
	})
	
	return &SimpleFSA{keys: keysCopy}
}

// Contains returns true if the key is in the set.
func (fsa *SimpleFSA) Contains(key []byte) bool {
	// Binary search
	left, right := 0, len(fsa.keys)
	for left < right {
		mid := (left + right) / 2
		cmp := bytes.Compare(fsa.keys[mid], key)
		if cmp < 0 {
			left = mid + 1
		} else if cmp > 0 {
			right = mid
		} else {
			return true
		}
	}
	return false
}

// Iterator returns an iterator over all keys in lexicographic order.
func (fsa *SimpleFSA) Iterator() FSAIterator {
	return &SimpleFSAIterator{
		fsa: fsa,
		pos: -1,
	}
}

// PrefixIterator returns an iterator over keys with the given prefix.
func (fsa *SimpleFSA) PrefixIterator(prefix []byte) FSAIterator {
	return &SimpleFSAIterator{
		fsa:    fsa,
		pos:    -1,
		prefix: prefix,
	}
}

// RangeIterator returns an iterator over keys in the given range [start, end).
func (fsa *SimpleFSA) RangeIterator(start, end []byte) FSAIterator {
	return &SimpleFSAIterator{
		fsa:   fsa,
		pos:   -1,
		start: start,
		end:   end,
	}
}

// Len returns the number of keys in the set.
func (fsa *SimpleFSA) Len() int {
	return len(fsa.keys)
}

// NumStates returns the total number of states.
func (fsa *SimpleFSA) NumStates() int {
	return len(fsa.keys) + 1 // Rough estimate
}

// SimpleFSAIterator implements FSAIterator for SimpleFSA.
type SimpleFSAIterator struct {
	fsa    *SimpleFSA
	pos    int
	prefix []byte
	start  []byte
	end    []byte
}

// Next advances the iterator and returns true if a value is available.
func (iter *SimpleFSAIterator) Next() bool {
	iter.pos++
	
	for iter.pos < len(iter.fsa.keys) {
		key := iter.fsa.keys[iter.pos]
		
		// Check prefix constraint
		if iter.prefix != nil && !bytes.HasPrefix(key, iter.prefix) {
			if bytes.Compare(key, iter.prefix) > 0 && !bytes.HasPrefix(iter.prefix, key) {
				return false
			}
			iter.pos++
			continue
		}
		
		// Check start constraint
		if iter.start != nil && bytes.Compare(key, iter.start) < 0 {
			iter.pos++
			continue
		}
		
		// Check end constraint
		if iter.end != nil && bytes.Compare(key, iter.end) >= 0 {
			return false
		}
		
		return true
	}
	
	return false
}

// Key returns the current key.
func (iter *SimpleFSAIterator) Key() []byte {
	if iter.pos >= 0 && iter.pos < len(iter.fsa.keys) {
		key := iter.fsa.keys[iter.pos]
		result := make([]byte, len(key))
		copy(result, key)
		return result
	}
	return nil
}

// Reset resets the iterator to the beginning.
func (iter *SimpleFSAIterator) Reset() {
	iter.pos = -1
}

// Seek positions the iterator at the first key >= target.
func (iter *SimpleFSAIterator) Seek(target []byte) bool {
	left, right := 0, len(iter.fsa.keys)
	for left < right {
		mid := (left + right) / 2
		if bytes.Compare(iter.fsa.keys[mid], target) < 0 {
			left = mid + 1
		} else {
			right = mid
		}
	}
	
	iter.pos = left - 1
	return iter.Next()
}

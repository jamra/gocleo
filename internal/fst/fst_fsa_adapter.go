package fst

import "bytes"

// FSTFSAAdapter adapts FST to implement the FSA interface for automata intersection
type FSTFSAAdapter struct {
	fst *FST
}

// NewFSTFSAAdapter creates an adapter to make FST implement FSA interface
func NewFSTFSAAdapter(fst *FST) *FSTFSAAdapter {
	return &FSTFSAAdapter{fst: fst}
}

// Contains checks if a key exists in the FST
func (adapter *FSTFSAAdapter) Contains(key []byte) bool {
	return adapter.fst.Contains(key)
}

// Iterator returns an iterator over all keys
func (adapter *FSTFSAAdapter) Iterator() FSAIterator {
	return &FSTFSAIteratorAdapter{iter: adapter.fst.Iterator()}
}

// PrefixIterator returns an iterator over keys with the given prefix
func (adapter *FSTFSAAdapter) PrefixIterator(prefix []byte) FSAIterator {
	return &FSTFSAIteratorAdapter{iter: adapter.fst.PrefixIterator(prefix)}
}

// RangeIterator returns an iterator over keys in the given range
func (adapter *FSTFSAAdapter) RangeIterator(start, end []byte) FSAIterator {
	return &FSTFSAIteratorAdapter{iter: adapter.fst.RangeIterator(start, end)}
}

// Len returns the number of keys
func (adapter *FSTFSAAdapter) Len() int {
	return adapter.fst.Size()
}

// NumStates returns an approximation of the number of states (simplified)
func (adapter *FSTFSAAdapter) NumStates() int {
	// For a simple FST implementation, approximate as number of unique prefixes
	// This is a rough estimate since we don't have access to internal state structure
	return adapter.fst.Size()
}

// FSTFSAIteratorAdapter adapts FST iterators to FSA iterator interface
type FSTFSAIteratorAdapter struct {
	iter interface{} // Could be *FSTIterator, *FSTPrefixIterator, or *FSTRangeIterator
	currentKey []byte
}

// Next advances the iterator and returns true if there's a next element
func (adapter *FSTFSAIteratorAdapter) Next() bool {
	switch iter := adapter.iter.(type) {
	case *FSTIterator:
		if iter.HasNext() {
			key, _ := iter.Next()
			adapter.currentKey = key
			return true
		}
	case *FSTPrefixIterator:
		if iter.HasNext() {
			key, _ := iter.Next()
			adapter.currentKey = key
			return true
		}
	case *FSTRangeIterator:
		if iter.HasNext() {
			key, _ := iter.Next()
			adapter.currentKey = key
			return true
		}
	}
	return false
}

// Key returns the current key
func (adapter *FSTFSAIteratorAdapter) Key() []byte {
	return adapter.currentKey
}

// Reset resets the iterator (simplified implementation)
func (adapter *FSTFSAIteratorAdapter) Reset() {
	// For simplicity, we'll just reset the current key
	// A full implementation would reset the underlying iterator
	adapter.currentKey = nil
}

// Seek seeks to the target key (simplified implementation)
func (adapter *FSTFSAIteratorAdapter) Seek(target []byte) bool {
	// For simplicity, we'll just iterate until we find the target
	// A full implementation would use binary search or FST-specific seeking
	for adapter.Next() {
		if bytes.Compare(adapter.currentKey, target) >= 0 {
			return true
		}
	}
	return false
}
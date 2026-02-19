package fst

import (
	"sort"
)

// SetOperationType represents the type of set operation
type SetOperationType int

const (
	UnionOp SetOperationType = iota
	IntersectionOp
	DifferenceOp
	SymmetricDifferenceOp
)

// SetOperation performs set operations between multiple FSAs
type SetOperation struct {
	fsas      []FSA
	operation SetOperationType
}

// NewSetOperation creates a new set operation
func NewSetOperation(operation SetOperationType, fsas ...FSA) *SetOperation {
	return &SetOperation{
		fsas:      fsas,
		operation: operation,
	}
}

// Execute performs the set operation and returns the result as a new FSA
func (so *SetOperation) Execute() (FSA, error) {
	if len(so.fsas) == 0 {
		return NewFSABuilder().Build()
	}
	
	if len(so.fsas) == 1 {
		return so.fsas[0], nil
	}
	
	var resultKeys []string
	
	switch so.operation {
	case UnionOp:
		resultKeys = so.executeUnion()
	case IntersectionOp:
		resultKeys = so.executeIntersection()
	case DifferenceOp:
		resultKeys = so.executeDifference()
	case SymmetricDifferenceOp:
		resultKeys = so.executeSymmetricDifference()
	}
	
	// Build result FSA
	builder := NewFSABuilder()
	for _, key := range resultKeys {
		builder.Add([]byte(key))
	}
	
	return builder.Build()
}

// executeUnion performs union operation
func (so *SetOperation) executeUnion() []string {
	keySet := make(map[string]bool)
	
	// Add all keys from all FSAs
	for _, fsa := range so.fsas {
		iter := fsa.Iterator()
		for iter.Next() {
			key := string(iter.Key())
			keySet[key] = true
		}
	}
	
	// Convert to sorted slice
	result := make([]string, 0, len(keySet))
	for key := range keySet {
		result = append(result, key)
	}
	
	sort.Strings(result)
	return result
}

// executeIntersection performs intersection operation
func (so *SetOperation) executeIntersection() []string {
	if len(so.fsas) == 0 {
		return []string{}
	}
	
	// Start with first FSA's keys
	candidates := make(map[string]bool)
	iter := so.fsas[0].Iterator()
	for iter.Next() {
		key := string(iter.Key())
		candidates[key] = true
	}
	
	// Check each candidate against all other FSAs
	for i := 1; i < len(so.fsas); i++ {
		fsa := so.fsas[i]
		newCandidates := make(map[string]bool)
		
		for candidate := range candidates {
			if fsa.Contains([]byte(candidate)) {
				newCandidates[candidate] = true
			}
		}
		
		candidates = newCandidates
		
		// Early termination if no candidates left
		if len(candidates) == 0 {
			break
		}
	}
	
	// Convert to sorted slice
	result := make([]string, 0, len(candidates))
	for key := range candidates {
		result = append(result, key)
	}
	
	sort.Strings(result)
	return result
}

// executeDifference performs difference operation (first FSA minus others)
func (so *SetOperation) executeDifference() []string {
	if len(so.fsas) == 0 {
		return []string{}
	}
	
	// Start with first FSA's keys
	result := make(map[string]bool)
	iter := so.fsas[0].Iterator()
	for iter.Next() {
		key := string(iter.Key())
		result[key] = true
	}
	
	// Remove keys that exist in any other FSA
	for i := 1; i < len(so.fsas); i++ {
		fsa := so.fsas[i]
		for key := range result {
			if fsa.Contains([]byte(key)) {
				delete(result, key)
			}
		}
	}
	
	// Convert to sorted slice
	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	
	sort.Strings(keys)
	return keys
}

// executeSymmetricDifference performs symmetric difference operation
func (so *SetOperation) executeSymmetricDifference() []string {
	if len(so.fsas) != 2 {
		// Symmetric difference typically works with 2 sets
		// For multiple sets, we'll do it pairwise
		current := so.fsas[0]
		for i := 1; i < len(so.fsas); i++ {
			symDiff := NewSetOperation(SymmetricDifferenceOp, current, so.fsas[i])
			result, _ := symDiff.Execute()
			current = result
		}
		
		var keys []string
		iter := current.Iterator()
		for iter.Next() {
			keys = append(keys, string(iter.Key()))
		}
		return keys
	}
	
	fsa1, fsa2 := so.fsas[0], so.fsas[1]
	result := make(map[string]bool)
	
	// Add keys from fsa1 that are not in fsa2
	iter1 := fsa1.Iterator()
	for iter1.Next() {
		key := string(iter1.Key())
		if !fsa2.Contains([]byte(key)) {
			result[key] = true
		}
	}
	
	// Add keys from fsa2 that are not in fsa1
	iter2 := fsa2.Iterator()
	for iter2.Next() {
		key := string(iter2.Key())
		if !fsa1.Contains([]byte(key)) {
			result[key] = true
		}
	}
	
	// Convert to sorted slice
	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	
	sort.Strings(keys)
	return keys
}

// Convenient methods for FSA

// Union returns the union of this FSA with others
func Union(fsa FSA, others ...FSA) (FSA, error) {
	fsas := append([]FSA{fsa}, others...)
	op := NewSetOperation(UnionOp, fsas...)
	return op.Execute()
}

// Intersection returns the intersection of this FSA with others
func Intersection(fsa FSA, others ...FSA) (FSA, error) {
	fsas := append([]FSA{fsa}, others...)
	op := NewSetOperation(IntersectionOp, fsas...)
	return op.Execute()
}

// Difference returns this FSA minus the others
func Difference(fsa FSA, others ...FSA) (FSA, error) {
	fsas := append([]FSA{fsa}, others...)
	op := NewSetOperation(DifferenceOp, fsas...)
	return op.Execute()
}

// SymmetricDifference returns the symmetric difference with another FSA
func SymmetricDifference(fsa FSA, other FSA) (FSA, error) {
	op := NewSetOperation(SymmetricDifferenceOp, fsa, other)
	return op.Execute()
}
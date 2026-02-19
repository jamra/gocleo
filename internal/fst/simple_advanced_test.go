package fst

import (
	"testing"
)

// Test just the basic fuzzy search functionality
func TestSimpleFuzzySearch(t *testing.T) {
	// Create FSA with test words
	builder := NewFSABuilder()
	words := []string{"cat", "car", "dog"}
	
	for _, word := range words {
		builder.Add([]byte(word))
	}
	
	fsa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}
	
	// Test fuzzy search - should find "car" when searching for "cat" with distance 1
	results := FuzzySearch(fsa, "cat", 1)
	
	// Should find "cat" (exact match) and "car" (1 edit distance)
	if len(results) < 1 {
		t.Errorf("Expected at least 1 result, got %d", len(results))
	}
	
	// Verify "cat" is in results
	found := false
	for _, result := range results {
		if result == "cat" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'cat' in results, got %v", results)
	}
}

// Test basic regex search functionality
func TestSimpleRegexSearch(t *testing.T) {
	// Create FSA with test words
	builder := NewFSABuilder()
	words := []string{"apple", "application", "banana"}
	
	for _, word := range words {
		builder.Add([]byte(word))
	}
	
	fsa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}
	
	// Test regex search
	results, err := RegexSearch(fsa, "app.*")
	if err != nil {
		t.Fatalf("RegexSearch failed: %v", err)
	}
	
	// Should find "apple" and "application"
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d: %v", len(results), results)
	}
}

// Test basic set operations
func TestSimpleSetOperations(t *testing.T) {
	// Create first FSA
	builder1 := NewFSABuilder()
	words1 := []string{"apple", "banana"}
	for _, word := range words1 {
		builder1.Add([]byte(word))
	}
	fsa1, _ := builder1.Build()
	
	// Create second FSA
	builder2 := NewFSABuilder()
	words2 := []string{"banana", "cherry"}
	for _, word := range words2 {
		builder2.Add([]byte(word))
	}
	fsa2, _ := builder2.Build()
	
	// Test Union
	union, err := Union(fsa1, fsa2)
	if err != nil {
		t.Fatalf("Union failed: %v", err)
	}
	
	if union.Len() != 3 { // apple, banana, cherry
		t.Errorf("Union size = %d, expected 3", union.Len())
	}
	
	// Test Intersection
	intersection, err := Intersection(fsa1, fsa2)
	if err != nil {
		t.Fatalf("Intersection failed: %v", err)
	}
	
	if intersection.Len() != 1 { // banana
		t.Errorf("Intersection size = %d, expected 1", intersection.Len())
	}
}
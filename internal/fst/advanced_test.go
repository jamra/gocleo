package fst

import (
	"fmt"
	"testing"
)

// Test Automaton functionality
func TestAutomaton(t *testing.T) {
	builder := NewAutomatonBuilder()
	
	// Build automaton for words: "cat", "car", "card"
	words := []string{"car", "card", "cat"}
	automaton := builder.BuildFromStrings(words)
	
	// Test acceptance
	testCases := []struct {
		input    string
		expected bool
	}{
		{"car", true},
		{"card", true},
		{"cat", true},
		{"ca", false},
		{"cars", false},
		{"dog", false},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := automaton.Accept([]byte(tc.input))
		if result != tc.expected {
			t.Errorf("Accept(%s) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

// Test Levenshtein automaton
func TestLevenshteinAutomaton(t *testing.T) {
	// Create FSA with test words
	builder := NewFSABuilder()
	words := []string{"cat", "car", "card", "care", "careful", "dog", "dogs"}
	
	for _, word := range words {
		builder.Add([]byte(word))
	}
	
	fsa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}
	
	// Test fuzzy search
	results := FuzzySearch(fsa, "car", 1)
	
	expected := []string{"car", "cat"}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}
	
	for i, result := range results {
		if i < len(expected) && result != expected[i] {
			t.Errorf("Expected result[%d] = %s, got %s", i, expected[i], result)
		}
	}
	
	// Test fuzzy search with distance 2
	results2 := FuzzySearch(fsa, "care", 2)
	if len(results2) < 2 {
		t.Errorf("Expected at least 2 results for distance 2, got %d", len(results2))
	}
}

// Test Set Operations
func TestSetOperations(t *testing.T) {
	// Create first FSA
	builder1 := NewFSABuilder()
	words1 := []string{"apple", "banana", "cherry"}
	for _, word := range words1 {
		builder1.Add([]byte(word))
	}
	fsa1, _ := builder1.Build()
	
	// Create second FSA
	builder2 := NewFSABuilder()
	words2 := []string{"banana", "cherry", "date"}
	for _, word := range words2 {
		builder2.Add([]byte(word))
	}
	fsa2, _ := builder2.Build()
	
	// Test Union
	union, err := Union(fsa1, fsa2)
	if err != nil {
		t.Fatalf("Union failed: %v", err)
	}
	
	expectedUnion := []string{"apple", "banana", "cherry", "date"}
	if union.Len() != len(expectedUnion) {
		t.Errorf("Union size = %d, expected %d", union.Len(), len(expectedUnion))
	}
	
	for _, word := range expectedUnion {
		if !union.Contains([]byte(word)) {
			t.Errorf("Union missing word: %s", word)
		}
	}
	
	// Test Intersection
	intersection, err := Intersection(fsa1, fsa2)
	if err != nil {
		t.Fatalf("Intersection failed: %v", err)
	}
	
	expectedIntersection := []string{"banana", "cherry"}
	if intersection.Len() != len(expectedIntersection) {
		t.Errorf("Intersection size = %d, expected %d", intersection.Len(), len(expectedIntersection))
	}
	
	for _, word := range expectedIntersection {
		if !intersection.Contains([]byte(word)) {
			t.Errorf("Intersection missing word: %s", word)
		}
	}
	
	// Test Difference
	difference, err := Difference(fsa1, fsa2)
	if err != nil {
		t.Fatalf("Difference failed: %v", err)
	}
	
	expectedDifference := []string{"apple"}
	if difference.Len() != len(expectedDifference) {
		t.Errorf("Difference size = %d, expected %d", difference.Len(), len(expectedDifference))
	}
	
	if !difference.Contains([]byte("apple")) {
		t.Errorf("Difference missing word: apple")
	}
}

// Test Regex functionality
func TestRegexSearch(t *testing.T) {
	// Create FSA with test words
	builder := NewFSABuilder()
	words := []string{"apple", "application", "apply", "banana", "band", "bandana"}
	
	for _, word := range words {
		builder.Add([]byte(word))
	}
	
	fsa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}
	
	// Test basic regex
	results, err := RegexSearch(fsa, "app.*")
	if err != nil {
		t.Fatalf("RegexSearch failed: %v", err)
	}
	
	expected := []string{"apple", "application", "apply"}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}
	
	// Test prefix + regex
	prefixResults, err := PrefixRegexSearch(fsa, "ban", "ban.*a")
	if err != nil {
		t.Fatalf("PrefixRegexSearch failed: %v", err)
	}
	
	expectedPrefix := []string{"banana", "bandana"}
	if len(prefixResults) != len(expectedPrefix) {
		t.Errorf("Expected %d prefix results, got %d", len(expectedPrefix), len(prefixResults))
	}
}

// Test Complex Queries
func TestComplexQueries(t *testing.T) {
	// Create FSA with test data
	builder := NewFSABuilder()
	words := []string{"apple", "application", "apply", "banana", "band", "bandana", "car", "card", "care"}
	
	for _, word := range words {
		builder.Add([]byte(word))
	}
	
	fsa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}
	
	query := NewComplexQuery(fsa)
	
	// Test prefix + regex combination
	options := QueryOptions{
		Prefix:       "app",
		RegexPattern: "app.*e",
		Limit:        10,
	}
	
	result, err := query.Execute(options)
	if err != nil {
		t.Fatalf("Complex query failed: %v", err)
	}
	
	expected := []string{"apple"}
	if len(result.Keys) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(result.Keys))
	}
	
	if result.Count != len(result.Keys) {
		t.Errorf("Result count mismatch: %d vs %d", result.Count, len(result.Keys))
	}
}

// Test FST functionality
func TestFST(t *testing.T) {
	builder := NewFSTBuilder()
	
	// Add key-value pairs
	testData := []struct {
		key   string
		value uint64
	}{
		{"apple", 1},
		{"banana", 2},
		{"cherry", 3},
	}
	
	for _, item := range testData {
		err := builder.Add([]byte(item.key), item.value)
		if err != nil {
			t.Fatalf("Failed to add %s: %v", item.key, err)
		}
	}
	
	fst, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FST: %v", err)
	}
	
	// Test retrieval
	for _, item := range testData {
		value, exists := fst.Get([]byte(item.key))
		if !exists {
			t.Errorf("Key %s not found", item.key)
		}
		if value != item.value {
			t.Errorf("Key %s: expected value %d, got %d", item.key, item.value, value)
		}
	}
	
	// Test non-existent key
	_, exists := fst.Get([]byte("nonexistent"))
	if exists {
		t.Errorf("Non-existent key should not be found")
	}
	
	// Test size
	if fst.Size() != len(testData) {
		t.Errorf("Size = %d, expected %d", fst.Size(), len(testData))
	}
}

// Test FST Iterators
func TestFSTIterators(t *testing.T) {
	builder := NewFSTBuilder()
	
	testData := []struct {
		key   string
		value uint64
	}{
		{"apple", 1},
		{"banana", 2},
		{"cherry", 3},
		{"date", 4},
	}
	
	for _, item := range testData {
		builder.Add([]byte(item.key), item.value)
	}
	
	fst, _ := builder.Build()
	
	// Test full iterator
	iter := fst.Iterator()
	count := 0
	for iter.HasNext() {
		key, value := iter.Next()
		if count < len(testData) {
			expectedKey := testData[count].key
			expectedValue := testData[count].value
			
			if string(key) != expectedKey {
				t.Errorf("Iterator key[%d] = %s, expected %s", count, string(key), expectedKey)
			}
			if value != expectedValue {
				t.Errorf("Iterator value[%d] = %d, expected %d", count, value, expectedValue)
			}
		}
		count++
	}
	
	if count != len(testData) {
		t.Errorf("Iterator count = %d, expected %d", count, len(testData))
	}
}

// Test FST Set Operations
func TestFSTSetOperations(t *testing.T) {
	// Create first FST
	builder1 := NewFSTBuilder()
	builder1.Add([]byte("apple"), 1)
	builder1.Add([]byte("banana"), 2)
	fst1, _ := builder1.Build()
	
	// Create second FST
	builder2 := NewFSTBuilder()
	builder2.Add([]byte("banana"), 20) // Different value
	builder2.Add([]byte("cherry"), 3)
	fst2, _ := builder2.Build()
	
	// Test union
	union, err := FSTUnion(fst1, fst2)
	if err != nil {
		t.Fatalf("FST Union failed: %v", err)
	}
	
	if union.Size() != 3 {
		t.Errorf("Union size = %d, expected 3", union.Size())
	}
	
	// Check that first FST's value takes precedence for "banana"
	value, exists := union.Get([]byte("banana"))
	if !exists || value != 2 {
		t.Errorf("Union banana value = %d, expected 2", value)
	}
	
	// Test intersection
	intersection, err := FSTIntersection(fst1, fst2)
	if err != nil {
		t.Fatalf("FST Intersection failed: %v", err)
	}
	
	if intersection.Size() != 1 {
		t.Errorf("Intersection size = %d, expected 1", intersection.Size())
	}
	
	// Should only contain "banana" with first FST's value
	value, exists = intersection.Get([]byte("banana"))
	if !exists || value != 2 {
		t.Errorf("Intersection banana value = %d, expected 2", value)
	}
}

// Benchmark tests
func BenchmarkFuzzySearch(b *testing.B) {
	// Create large FSA
	builder := NewFSABuilder()
	for i := 0; i < 1000; i++ {
		builder.Add([]byte(fmt.Sprintf("word%04d", i)))
	}
	fsa, _ := builder.Build()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FuzzySearch(fsa, "word0500", 2)
	}
}

func BenchmarkRegexSearch(b *testing.B) {
	// Create large FSA
	builder := NewFSABuilder()
	for i := 0; i < 1000; i++ {
		builder.Add([]byte(fmt.Sprintf("item%04d", i)))
	}
	fsa, _ := builder.Build()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegexSearch(fsa, "item.*5.*")
	}
}

func BenchmarkSetUnion(b *testing.B) {
	// Create two FSAs
	builder1 := NewFSABuilder()
	for i := 0; i < 500; i++ {
		builder1.Add([]byte(fmt.Sprintf("word%04d", i)))
	}
	fsa1, _ := builder1.Build()
	
	builder2 := NewFSABuilder()
	for i := 250; i < 750; i++ {
		builder2.Add([]byte(fmt.Sprintf("word%04d", i)))
	}
	fsa2, _ := builder2.Build()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Union(fsa1, fsa2)
	}
}
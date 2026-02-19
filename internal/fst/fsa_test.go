package fst

import (
	"fmt"
	"sort"
	"testing"
)

func TestSimpleFSABasicOperations(t *testing.T) {
	words := []string{
		"apple",
		"application",
		"apply",
		"banana",
		"band",
		"bandana",
		"cat",
		"category",
	}

	// Build FSA
	fsaBuilder := NewFSABuilder()
	for _, word := range words {
		err := fsaBuilder.Add([]byte(word))
		if err != nil {
			t.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fsa, err := fsaBuilder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}

	// Test Contains
	for _, word := range words {
		if !fsa.Contains([]byte(word)) {
			t.Errorf("FSA should contain '%s'", word)
		}
	}

	// Test non-existent words
	nonExistentWords := []string{
		"app", "appl", "applications", "ban", "bands", "categories", "dog",
	}

	for _, word := range nonExistentWords {
		if fsa.Contains([]byte(word)) {
			t.Errorf("FSA should not contain '%s'", word)
		}
	}

	// Test Len
	if fsa.Len() != len(words) {
		t.Errorf("Expected FSA length %d, got %d", len(words), fsa.Len())
	}
}

func TestSimpleFSAIterator(t *testing.T) {
	words := []string{"a", "aa", "aaa", "b", "bb", "c"}

	// Build FSA
	fsaBuilder := NewFSABuilder()
	for _, word := range words {
		err := fsaBuilder.Add([]byte(word))
		if err != nil {
			t.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fsa, err := fsaBuilder.Build()
	if err != nil {
		t.Fatalf("Failed to build FSA: %v", err)
	}

	// Test full iteration
	iter := fsa.Iterator()
	var results []string

	for iter.Next() {
		key := iter.Key()
		results = append(results, string(key))
	}

	// Verify results are in lexicographic order
	if !sort.StringsAreSorted(results) {
		t.Errorf("Iterator results are not sorted: %v", results)
	}

	// Verify all words are present
	if len(results) != len(words) {
		t.Errorf("Expected %d results, got %d", len(words), len(results))
	}
}

func TestSimpleFSAPrefixIterator(t *testing.T) {
	words := []string{
		"apple", "application", "apply", "banana", "band", "bandana", "cat", "category",
	}

	fsaBuilder := NewFSABuilder()
	for _, word := range words {
		fsaBuilder.Add([]byte(word))
	}

	fsa, _ := fsaBuilder.Build()

	// Test prefix "app"
	prefixIter := fsa.PrefixIterator([]byte("app"))
	var appResults []string

	for prefixIter.Next() {
		key := prefixIter.Key()
		appResults = append(appResults, string(key))
	}

	expectedAppWords := []string{"apple", "application", "apply"}
	if len(appResults) != len(expectedAppWords) {
		t.Errorf("Expected %d words with prefix 'app', got %d: %v",
			len(expectedAppWords), len(appResults), appResults)
	}
}

func TestFSABuilderSortedInput(t *testing.T) {
	fsaBuilder := NewFSABuilder()

	// Test adding unsorted keys should fail
	fsaBuilder.Add([]byte("zebra"))
	err := fsaBuilder.Add([]byte("apple"))
	if err == nil {
		t.Error("Expected error when adding unsorted key")
	}

	// Test duplicate keys should fail
	fsaBuilder.Reset()
	fsaBuilder.Add([]byte("apple"))
	err = fsaBuilder.Add([]byte("apple"))
	if err == nil {
		t.Error("Expected error when adding duplicate key")
	}
}

func BenchmarkSimpleFSABuild(b *testing.B) {
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%05d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fsaBuilder := NewFSABuilder()
		for _, word := range words {
			fsaBuilder.Add([]byte(word))
		}
		fsaBuilder.Build()
	}
}

func BenchmarkSimpleFSALookup(b *testing.B) {
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%05d", i)
	}

	fsaBuilder := NewFSABuilder()
	for _, word := range words {
		fsaBuilder.Add([]byte(word))
	}
	fsa, _ := fsaBuilder.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := words[i%len(words)]
		fsa.Contains([]byte(word))
	}
}

package main

import (
	"fmt"
	"log"
	
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== Advanced FSA/FST Features Demo ===\n")
	
	// Demo 1: Basic FSA with advanced search
	fmt.Println("1. Building FSA with sample data...")
	fsaBuilder := fst.NewFSABuilder()
	
	words := []string{
		"apple", "application", "apply", "approach",
		"banana", "band", "bandage", "bandana",
		"car", "card", "care", "careful", "careless",
		"dog", "dogma", "dogmatic", "dogs",
		"treat", "treatment", "tree", "trees",
	}
	
	for _, word := range words {
		err := fsaBuilder.Add([]byte(word))
		if err != nil {
			log.Fatalf("Failed to add word %s: %v", word, err)
		}
	}
	
	fsa, err := fsaBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build FSA: %v", err)
	}
	
	fmt.Printf("FSA built with %d words\n\n", fsa.Len())
	
	// Demo 2: Fuzzy Search
	fmt.Println("2. Fuzzy Search Demo:")
	pattern := "care"
	maxDistance := 2
	fuzzyResults := fst.FuzzySearch(fsa, pattern, maxDistance)
	
	fmt.Printf("Fuzzy search for '%s' (max distance %d):\n", pattern, maxDistance)
	for _, result := range fuzzyResults {
		fmt.Printf("  - %s\n", result)
	}
	fmt.Println()
	
	// Demo 3: Regex Search
	fmt.Println("3. Regex Search Demo:")
	regexPattern := "app.*"
	regexResults, err := fst.RegexSearch(fsa, regexPattern)
	if err != nil {
		log.Fatalf("Regex search failed: %v", err)
	}
	
	fmt.Printf("Regex search for pattern '%s':\n", regexPattern)
	for _, result := range regexResults {
		fmt.Printf("  - %s\n", result)
	}
	fmt.Println()
	
	// Demo 4: Prefix + Regex combination
	fmt.Println("4. Prefix + Regex Search:")
	prefix := "ban"
	prefixRegex := "ban.*a"
	prefixResults, err := fst.PrefixRegexSearch(fsa, prefix, prefixRegex)
	if err != nil {
		log.Fatalf("Prefix regex search failed: %v", err)
	}
	
	fmt.Printf("Words starting with '%s' matching pattern '%s':\n", prefix, prefixRegex)
	for _, result := range prefixResults {
		fmt.Printf("  - %s\n", result)
	}
	fmt.Println()
	
	// Demo 5: Set Operations
	fmt.Println("5. Set Operations Demo:")
	
	// Create second FSA for set operations
	fsaBuilder2 := fst.NewFSABuilder()
	words2 := []string{"apple", "apricot", "banana", "blueberry", "car", "carrot"}
	
	for _, word := range words2 {
		fsaBuilder2.Add([]byte(word))
	}
	fsa2, _ := fsaBuilder2.Build()
	
	// Union
	union, err := fst.Union(fsa, fsa2)
	if err != nil {
		log.Fatalf("Union failed: %v", err)
	}
	
	fmt.Printf("Union of FSA1 (%d words) and FSA2 (%d words) = %d words\n", 
		fsa.Len(), fsa2.Len(), union.Len())
	
	// Intersection
	intersection, err := fst.Intersection(fsa, fsa2)
	if err != nil {
		log.Fatalf("Intersection failed: %v", err)
	}
	
	fmt.Printf("Intersection = %d words: ", intersection.Len())
	iter := intersection.Iterator()
	for iter.Next() {
		word := iter.Key()
		fmt.Printf("%s ", string(word))
	}
	fmt.Println()
	
	// Difference
	difference, err := fst.Difference(fsa, fsa2)
	if err != nil {
		log.Fatalf("Difference failed: %v", err)
	}
	
	fmt.Printf("Difference (FSA1 - FSA2) = %d words\n", difference.Len())
	fmt.Println()
	
	// Demo 6: Complex Query
	fmt.Println("6. Complex Query Demo:")
	query := fst.NewComplexQuery(fsa)
	
	options := fst.QueryOptions{
		Prefix:           "car",
		RegexPattern:     "car.*",
		FuzzyPattern:     "care",
		FuzzyMaxDistance: 1,
		Limit:           5,
	}
	
	result, err := query.Execute(options)
	if err != nil {
		log.Fatalf("Complex query failed: %v", err)
	}
	
	fmt.Printf("Complex query results (%d found):\n", result.Count)
	for _, key := range result.Keys {
		fmt.Printf("  - %s\n", key)
	}
	fmt.Println()
	
	// Demo 7: FST (Key-Value) Operations
	fmt.Println("7. FST (Key-Value) Demo:")
	
	fstBuilder := fst.NewFSTBuilder()
	keyValues := []struct {
		key   string
		value uint64
	}{
		{"apple", 100},
		{"banana", 85},
		{"cherry", 120},
		{"date", 90},
		{"elderberry", 75},
	}
	
	for _, kv := range keyValues {
		err := fstBuilder.Add([]byte(kv.key), kv.value)
		if err != nil {
			log.Fatalf("Failed to add key-value %s:%d: %v", kv.key, kv.value, err)
		}
	}
	
	fstMap, err := fstBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build FST: %v", err)
	}
	
	fmt.Printf("FST built with %d key-value pairs:\n", fstMap.Size())
	
	// Iterate through all key-value pairs
	fstIter := fstMap.Iterator()
	for fstIter.HasNext() {
		key, value := fstIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
	}
	fmt.Println()
	
	// Range query on FST
	fmt.Println("Range query (b to d):")
	rangeIter := fstMap.RangeIterator([]byte("b"), []byte("d"))
	for rangeIter.HasNext() {
		key, value := rangeIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
	}
	fmt.Println()
	
	// Demo 8: FST Set Operations
	fmt.Println("8. FST Set Operations:")
	
	// Create second FST
	fstBuilder2 := fst.NewFSTBuilder()
	keyValues2 := []struct {
		key   string
		value uint64
	}{
		{"banana", 95}, // Different value
		{"cherry", 120}, // Same value
		{"fig", 110},
		{"grape", 130},
	}
	
	for _, kv := range keyValues2 {
		fstBuilder2.Add([]byte(kv.key), kv.value)
	}
	fst2Map, _ := fstBuilder2.Build()
	
	// FST Union
	fstUnion, err := fst.FSTUnion(fstMap, fst2Map)
	if err != nil {
		log.Fatalf("FST union failed: %v", err)
	}
	
	fmt.Printf("FST Union (%d + %d = %d pairs):\n", 
		fstMap.Size(), fst2Map.Size(), fstUnion.Size())
	
	unionIter := fstUnion.Iterator()
	for unionIter.HasNext() {
		key, value := unionIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
	}
	fmt.Println()
	
	// FST Intersection
	fstIntersection, err := fst.FSTIntersection(fstMap, fst2Map)
	if err != nil {
		log.Fatalf("FST intersection failed: %v", err)
	}
	
	fmt.Printf("FST Intersection (%d pairs):\n", fstIntersection.Size())
	intersectionIter := fstIntersection.Iterator()
	for intersectionIter.HasNext() {
		key, value := intersectionIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
	}
	fmt.Println()
	
	// Demo 9: Performance Summary
	fmt.Println("9. Performance Summary:")
	fmt.Printf("  - Basic lookup: O(log n) binary search\n")
	fmt.Printf("  - Fuzzy search: O(m * n * d) where m=pattern, n=keys, d=distance\n")
	fmt.Printf("  - Regex search: O(n * m) where n=keys, m=pattern complexity\n")
	fmt.Printf("  - Set operations: O(n + m) where n,m are FSA sizes\n")
	fmt.Printf("  - Memory usage: ~8 bytes per key overhead\n")
	
	fmt.Println("\n=== Demo Complete! ===")
}
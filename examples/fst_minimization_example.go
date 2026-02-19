package main

import (
	"fmt"
	"log"
	
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== FST Minimization Algorithm Demo ===")
	fmt.Println()
	
	// Demonstrate basic minimization
	demonstrateBasicMinimization()
	
	// Demonstrate comparison with simple FST
	demonstrateComparison()
}

func demonstrateBasicMinimization() {
	fmt.Println("1. Basic Minimization Example")
	fmt.Println("-----------------------------")
	
	// Create a minimizing builder
	builder := fst.NewMinimizingBuilder()
	
	// Add some related words that should share states
	words := []struct {
		word  string
		value uint64
	}{
		{"car", 1},
		{"card", 2},
		{"care", 3},
		{"careful", 4},
		{"cat", 5},
		{"catch", 6},
		{"dog", 7},
		{"dogs", 8},
	}
	
	fmt.Printf("Adding %d words to FST...\n", len(words))
	for _, w := range words {
		err := builder.Add([]byte(w.word), w.value)
		if err != nil {
			log.Fatalf("Error adding word %s: %v", w.word, err)
		}
	}
	
	// Build the minimized FST
	minFST, err := builder.Build()
	if err != nil {
		log.Fatalf("Error building FST: %v", err)
	}
	
	fmt.Printf("✓ Built minimized FST with %d states\n", minFST.NumStates())
	fmt.Printf("✓ Estimated memory usage: %d bytes\n", minFST.EstimateMemoryUsage())
	
	// Test lookups
	fmt.Println("\nTesting lookups:")
	for _, w := range words {
		value, found := minFST.Get([]byte(w.word))
		if found {
			fmt.Printf("  ✓ %s -> %d\n", w.word, value)
		} else {
			fmt.Printf("  ✗ %s not found\n", w.word)
		}
	}
	
	// Test non-existent words
	fmt.Println("\nTesting non-existent words:")
	nonExistent := []string{"ca", "cars", "cares", "do"}
	for _, word := range nonExistent {
		if minFST.Contains([]byte(word)) {
			fmt.Printf("  ✗ %s found (shouldn't exist)\n", word)
		} else {
			fmt.Printf("  ✓ %s correctly not found\n", word)
		}
	}
	
	fmt.Println()
}

func demonstrateComparison() {
	fmt.Println("2. Comparison with Simple FST")
	fmt.Println("-----------------------------")
	
	// Test data with potential for sharing
	testData := []struct {
		key   string
		value uint64
	}{
		{"apple", 1},
		{"application", 2},
		{"apply", 3},
		{"banana", 4},
		{"band", 5},
		{"bandana", 6},
		{"cat", 7},
		{"catch", 8},
		{"dog", 9},
		{"doggy", 10},
	}
	
	// Build with minimizing builder
	fmt.Printf("Building minimized FST with %d entries...\n", len(testData))
	minBuilder := fst.NewMinimizingBuilder()
	for _, item := range testData {
		err := minBuilder.Add([]byte(item.key), item.value)
		if err != nil {
			log.Fatalf("Error adding to minimized builder: %v", err)
		}
	}
	
	minFST, err := minBuilder.Build()
	if err != nil {
		log.Fatalf("Error building minimized FST: %v", err)
	}
	
	// Build with simple builder
	fmt.Printf("Building simple FST for comparison...\n")
	simpleBuilder := fst.NewFSTBuilder()
	for _, item := range testData {
		err := simpleBuilder.Add([]byte(item.key), item.value)
		if err != nil {
			log.Fatalf("Error adding to simple builder: %v", err)
		}
	}
	
	simpleFST, err := simpleBuilder.Build()
	if err != nil {
		log.Fatalf("Error building simple FST: %v", err)
	}
	
	// Compare results
	fmt.Printf("\n📊 Comparison Results:\n")
	fmt.Printf("  Minimized FST: %d 'states', ~%d bytes\n", 
		minFST.NumStates(), minFST.EstimateMemoryUsage())
	fmt.Printf("  Simple FST: %d entries, array-based storage\n", simpleFST.Size())
	
	// Verify both give same results
	fmt.Printf("\n🔍 Verification (testing %d lookups)...\n", len(testData))
	allMatch := true
	for _, item := range testData {
		minVal, minFound := minFST.Get([]byte(item.key))
		simpleVal, simpleFound := simpleFST.Get([]byte(item.key))
		
		if !minFound || !simpleFound || minVal != simpleVal {
			fmt.Printf("  ✗ Mismatch for %s\n", item.key)
			allMatch = false
		}
	}
	
	if allMatch {
		fmt.Printf("  ✓ All lookups match between minimized and simple FST\n")
	}
	
	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nNote: This is a simplified implementation. A full minimization")
	fmt.Println("algorithm would use sophisticated state sharing techniques to")
	fmt.Println("dramatically reduce the number of states compared to naive approaches.")
}

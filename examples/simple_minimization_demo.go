package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== Simple FSA Demonstration ===\n")

	words := []string{"test", "testing", "best", "rest"}
	sort.Strings(words)

	fmt.Printf("Input words: %v\n\n", words)

	// Build FSA using the available builder
	builder := fst.NewFSABuilder()
	for _, word := range words {
		err := builder.Add([]byte(word))
		if err != nil {
			log.Fatalf("Error adding word: %v", err)
		}
	}

	fsa, err := builder.Build()
	if err != nil {
		log.Fatalf("Error building FSA: %v", err)
	}

	// Test the FSA
	fmt.Println("Testing FSA:")
	testWords := []string{"test", "testing", "xyz", "tes", "best"}
	
	for _, word := range testWords {
		result := fsa.Contains([]byte(word))
		fmt.Printf("  Contains '%s': %v\n", word, result)
	}

	// Test prefix iteration
	fmt.Println("\nPrefix 'test*' results:")
	iter := fsa.PrefixIterator([]byte("test"))
	for iter.Next() {
		fmt.Printf("  - %s\n", string(iter.Key()))
	}

	fmt.Println("\n✅ FSA demonstration complete!")
	fmt.Printf("Estimated size: %d bytes\n", builder.EstimatedSize())
}

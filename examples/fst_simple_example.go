package main

import (
	"fmt"
	"log"
	
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("FSA (Finite State Acceptor) Example")
	fmt.Println("===================================")
	
	// Sample dictionary words (must be in lexicographic order)
	words := []string{
		"apple",
		"application",
		"apply",
		"banana",
		"band",
		"bandana", 
		"cat",
		"category",
		"dog",
		"doghouse",
		"elephant",
		"example",
		"zebra",
	}
	
	// Build FSA
	fmt.Printf("Building FSA with %d words...\n", len(words))
	fsaBuilder := fst.NewFSABuilder()
	
	for _, word := range words {
		err := fsaBuilder.Add([]byte(word))
		if err != nil {
			log.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}
	
	fsa, err := fsaBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build FSA: %v", err)
	}
	
	fmt.Printf("FSA built successfully with %d words!\n\n", fsa.Len())
	
	// Test membership
	fmt.Println("1. Testing membership (Contains):")
	testWords := []string{"apple", "app", "application", "cats", "elephant", "zebra"}
	
	for _, word := range testWords {
		contains := fsa.Contains([]byte(word))
		fmt.Printf("   '%s': %v\n", word, contains)
	}
	
	// Test full iteration
	fmt.Println("\n2. Full iteration (all words in lexicographic order):")
	iter := fsa.Iterator()
	count := 0
	
	for iter.Next() {
		key := iter.Key()
		fmt.Printf("   %d. %s\n", count+1, string(key))
		count++
	}
	
	// Test prefix iteration
	fmt.Println("\n3. Prefix iteration (words starting with 'app'):")
	prefixIter := fsa.PrefixIterator([]byte("app"))
	count = 0
	
	for prefixIter.Next() {
		key := prefixIter.Key()
		fmt.Printf("   %d. %s\n", count+1, string(key))
		count++
	}
	
	// Test range iteration
	fmt.Println("\n4. Range iteration (words from 'banana' to 'dog'):")
	rangeIter := fsa.RangeIterator([]byte("banana"), []byte("dog"))
	count = 0
	
	for rangeIter.Next() {
		key := rangeIter.Key()
		fmt.Printf("   %d. %s\n", count+1, string(key))
		count++
	}
	
	// Display statistics
	fmt.Println("\n5. FSA Statistics:")
	fmt.Printf("   Total words: %d\n", fsa.Len())
	fmt.Printf("   Total states: %d\n", fsa.NumStates())
	fmt.Printf("   Estimated memory usage: %d bytes\n", fsaBuilder.EstimatedSize())
	
	fmt.Println("\nExample completed successfully!")
}

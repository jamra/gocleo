/*
 * Copyright (c) 2011 jamra.source@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

// Package main demonstrates basic FST (Finite State Transducer) usage
package main

import (
	"fmt"
	"log"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== FST Basic Example ===")

	// Create a new FST builder
	builder := fst.NewFSTBuilder()

	// Sample data - words with their document IDs
	words := []struct {
		word string
		id   uint64
	}{
		{"apple", 1},
		{"application", 2},
		{"apply", 3},
		{"banana", 4},
		{"band", 5},
		{"bandana", 6},
		{"car", 7},
		{"card", 8},
		{"care", 9},
		{"careful", 10},
	}

	// Add words to FST (they must be added in lexicographic order)
	fmt.Println("\nBuilding FST with words:")
	for _, word := range words {
		fmt.Printf("  Adding: %s -> %d\n", word.word, word.id)
		err := builder.Add([]byte(word.word), word.id)
		if err != nil {
			log.Fatalf("Error adding word '%s': %v", word.word, err)
		}
	}

	// Build the FST
	dictionary, err := builder.Build()
	if err != nil {
		log.Fatalf("Error building FST: %v", err)
	}

	fmt.Printf("\nFST built successfully! Size: %d entries\n", dictionary.Size())

	// Test exact lookups
	fmt.Println("\n=== Exact Lookups ===")
	testKeys := []string{"apple", "car", "banana", "xyz", "application"}

	for _, key := range testKeys {
		if value, found := dictionary.Get([]byte(key)); found {
			fmt.Printf("  '%s' -> %d ✓\n", key, value)
		} else {
			fmt.Printf("  '%s' -> not found ✗\n", key)
		}
	}

	// Test prefix search
	fmt.Println("\n=== Prefix Search ===")
	prefixes := []string{"app", "ban", "car"}

	for _, prefix := range prefixes {
		fmt.Printf("\nMatches for prefix '%s':\n", prefix)
		iter := dictionary.PrefixIterator([]byte(prefix))

		count := 0
		for iter.HasNext() {
			key, value := iter.Next()
			fmt.Printf("  %s -> %d\n", string(key), value)
			count++
		}

		if count == 0 {
			fmt.Printf("  No matches found\n")
		}
	}

	// Test range iteration
	fmt.Println("\n=== Range Iteration ===")
	fmt.Println("Words from 'app' to 'car':")
	rangeIter := dictionary.RangeIterator([]byte("app"), []byte("car"))

	count := 0
	for rangeIter.HasNext() {
		key, value := rangeIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
		count++
	}

	if count == 0 {
		fmt.Printf("  No matches found in range\n")
	}

	// Test full iteration
	fmt.Println("\n=== Full Dictionary ===")
	fullIter := dictionary.Iterator()

	fmt.Printf("All entries in the FST:\n")
	for fullIter.HasNext() {
		key, value := fullIter.Next()
		fmt.Printf("  %s -> %d\n", string(key), value)
	}

	fmt.Println("\n=== FST Basic Example Complete ===")
}

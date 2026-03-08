// Package main demonstrates true FST-Regex automata intersection
package main

import (
	"fmt"
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== Simple FST Automata Intersection Test ===")

	// Simple test data - manually sorted
	testData := []string{
		"app", "apple", "application",
		"banana", "band",
		"cat", "catch",
	}

	// Build FST
	builder := fst.NewFSTBuilder()
	for i, word := range testData {
		err := builder.Add([]byte(word), uint64(i))
		if err != nil {
			fmt.Printf("Failed to add %s: %v\n", word, err)
			return
		}
	}

	fstStructure, err := builder.Build()
	if err != nil {
		fmt.Printf("Failed to build FST: %v\n", err)
		return
	}

	// Create FSA adapter
	fsa := fst.NewFSTFSAAdapter(fstStructure)
	fmt.Printf("✅ FST built with %d keys\n", fsa.Len())

	// Test simple intersection
	pattern := "app.*"
	fmt.Printf("Testing pattern: %s\n", pattern)
	
	regexAutomaton, err := fst.NewTrueRegexAutomaton(pattern)
	if err != nil {
		fmt.Printf("Failed to create regex automaton: %v\n", err)
		return
	}

	results, err := regexAutomaton.IntersectWithFST(fsa)
	if err != nil {
		fmt.Printf("Intersection failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d results: %v\n", len(results), results)
	fmt.Println("✅ True automata intersection working!")
}
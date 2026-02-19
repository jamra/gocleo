package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("=== FSA Minimization Demonstration ===\n")

	words := []string{"test", "testing", "tester", "best", "besting", "bester", "rest", "resting"}
	sort.Strings(words)

	fmt.Printf("Input words: %v\n\n", words)

	// Compare different FSA implementations
	implementations := []struct {
		name    string
		options fst.FSAOptions
	}{
		{"Simple FSA", fst.FSAOptions{EnableAutomaton: false, EnableMinimization: false}},
		{"Automaton (unminimized)", fst.FSAOptions{EnableAutomaton: true, EnableMinimization: false}},
		{"Automaton (minimized)", fst.FSAOptions{EnableAutomaton: true, EnableMinimization: true}},
	}

	for _, impl := range implementations {
		fmt.Printf("%s:\n", impl.name)
		
		builder := fst.NewFSABuilderWithOptions(impl.options)
		for _, word := range words {
			builder.Add([]byte(word))
		}

		start := time.Now()
		fsa, err := builder.Build()
		buildTime := time.Since(start)
		
		if err != nil {
			log.Printf("Error building %s: %v", impl.name, err)
			continue
		}

		// Test lookups
		testWords := []string{"test", "xyz", "testing"}
		fmt.Printf("  Build time: %v\n", buildTime)
		fmt.Printf("  Test results: ")
		for _, word := range testWords {
			result := fsa.Contains([]byte(word))
			fmt.Printf("%s=%v ", word, result)
		}
		fmt.Printf("\n  Estimated size: %d bytes\n\n", builder.EstimatedSize())
	}

	fmt.Println("✅ Minimization demonstration complete!")
}

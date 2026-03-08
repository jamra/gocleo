package main

import (
	"fmt"
	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	// Build simple FST
	documents := []string{"apple", "application", "banana"}
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		panic(err)
	}

	// Create search engine
	searchEngine := fst.NewSearchEngine(fstIndex, documents, nil)

	// Test if methods exist
	fmt.Printf("FST has %d entries\n", fstIndex.Size())

	// Test intersection
	results, err := searchEngine.IntersectionRegexSearch("app.*")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Intersection results: %d\n", len(results))
	}

	// Test naive
	results, err = searchEngine.RegexSearch("app.*")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Naive results: %d\n", len(results))
	}
}
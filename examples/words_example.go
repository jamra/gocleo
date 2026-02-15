// Package main demonstrates creating a Cleo search from a word list.
package main

import (
	"fmt"
	"log"

	"github.com/jamra/gocleo/pkg/cleo"
)

func main() {
	// Create a search client from a list of words (useful for programmatic use)
	words := []string{
		"apple", "application", "apply", "approach",
		"banana", "band", "bank", "basketball",
		"computer", "compute", "company", "complete",
		"data", "database", "date", "development",
		"elephant", "email", "employee", "engineering",
	}

	config := &cleo.Config{
		ScoringFunction: cleo.PrefixScore, // Prioritize prefix matches
		MaxResults:      5,                // Limit to top 5 results
		MinScore:        0.1,              // Filter low-relevance results
	}

	client, err := cleo.NewFromWords(words, config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Test various queries
	queries := []string{"app", "comp", "dat", "ele", "ban"}

	fmt.Println("=== Word List Search Examples ===")
	for _, query := range queries {
		results, err := client.Search(query)
		if err != nil {
			log.Printf("Error searching for '%s': %v", query, err)
			continue
		}

		fmt.Printf("Query: '%s'\n", query)
		if len(results) == 0 {
			fmt.Println("  No results found")
		} else {
			for i, result := range results {
				fmt.Printf("  %d. %s (score: %.4f)\n", i+1, result.Word, result.Score)
			}
		}
		fmt.Println()
	}

	// Show how different scoring functions work
	fmt.Println("=== Scoring Function Comparison ===")
	query := "comp"

	scoringFunctions := map[string]cleo.ScoringFunction{
		"Default": cleo.DefaultScore,
		"Prefix":  cleo.PrefixScore,
		"Exact":   cleo.ExactScore,
		"Fuzzy":   cleo.FuzzyScore,
	}

	for name, scoringFunc := range scoringFunctions {
		client.SetScoringFunction(scoringFunc)
		results, _ := client.Search(query)
		
		fmt.Printf("%s scoring for '%s':\n", name, query)
		for i, result := range results {
			if i >= 3 { // Top 3
				break
			}
			fmt.Printf("  %d. %s (%.4f)\n", i+1, result.Word, result.Score)
		}
		fmt.Println()
	}
}

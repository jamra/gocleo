// Package main demonstrates the new Cleo API design.
package main

import (
	"fmt"
	"log"

	"github.com/jamra/gocleo/pkg/cleo"
	httpapi "github.com/jamra/gocleo/api/http"
)

func main() {
	// Example 1: Create a search client from a corpus file
	config := cleo.DefaultConfig()
	config.MaxResults = 10
	config.MinScore = 0.1

	client, err := cleo.New("./w1_fixed.txt", config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 2: Perform searches programmatically
	fmt.Println("=== Programmatic Search Examples ===")
	
	queries := []string{"pizza", "hello", "computer", "xyz"}
	for _, query := range queries {
		results, err := client.Search(query)
		if err != nil {
			log.Printf("Search error for '%s': %v", query, err)
			continue
		}
		
		fmt.Printf("Query: '%s' - %d results\n", query, len(results))
		for i, result := range results {
			if i >= 3 { // Show top 3
				break
			}
			fmt.Printf("  %d. %s (score: %.4f)\n", i+1, result.Word, result.Score)
		}
		fmt.Println()
	}

	// Example 3: Try different scoring functions
	fmt.Println("=== Different Scoring Functions ===")
	client.SetScoringFunction(cleo.PrefixScore)
	results, _ := client.Search("comp")
	fmt.Println("PrefixScore for 'comp':")
	for i, result := range results[:min(3, len(results))] {
		fmt.Printf("  %d. %s (score: %.4f)\n", i+1, result.Word, result.Score)
	}

	client.SetScoringFunction(cleo.FuzzyScore) 
	results, _ = client.Search("comp")
	fmt.Println("FuzzyScore for 'comp':")
	for i, result := range results[:min(3, len(results))] {
		fmt.Printf("  %d. %s (score: %.4f)\n", i+1, result.Word, result.Score)
	}

	// Example 4: Show index statistics
	fmt.Println("\n=== Index Statistics ===")
	stats := client.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Example 5: Start HTTP server (optional)
	fmt.Println("\n=== Starting HTTP Server ===")
	fmt.Println("Starting server on :8080")
	fmt.Println("Try: curl 'http://localhost:8080/search?q=pizza'")
	
	if err := httpapi.ListenAndServe("8080", client); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

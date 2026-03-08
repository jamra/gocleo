package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jamra/gocleo/internal/fst"
	"github.com/jamra/gocleo/internal/scoring"
)

func main() {
	fmt.Println("=== Simple FST Regex Example ===")

	// Sample documents
	documents := []string{
		"apple pie recipe",
		"banana bread instructions", 
		"cherry cake tutorial",
		"application development",
		"programming tutorial",
		"debugging techniques",
		"baking cookies delicious",
		"making bread fresh",
		"cooking pasta italian",
	}

	// Build FST from documents
	fmt.Println("Building FST from documents...")
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		log.Fatalf("Failed to build FST: %v", err)
	}

	// Create search engine
	searchEngine := fst.NewSearchEngine(fstIndex, documents, scoring.DefaultScore)
	fmt.Printf("FST built successfully with %d documents\n\n", len(documents))

	// Simple regex demonstration using basic FST functionality
	fmt.Println("=== Manual Regex Search (using FST + Go regex) ===")
	
	regexPatterns := []string{
		"app.*",      // Words starting with "app"
		".*ing$",     // Words ending with "ing"
		".*recipe.*", // Words containing "recipe"
	}

	for _, pattern := range regexPatterns {
		fmt.Printf("\nPattern: %s\n", pattern)
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Error compiling regex: %v\n", err)
			continue
		}

		matches := []string{}
		
		// Get all keys from FST and test against regex
		iter := fstIndex.Iterator()
		for iter.HasNext() {
			key, docID := iter.Next()
			keyStr := string(key)
			
			if regex.MatchString(keyStr) {
				if int(docID) < len(documents) && docID > 0 {
					docContent := documents[docID-1]
					matches = append(matches, docContent)
				}
			}
		}

		if len(matches) > 0 {
			fmt.Printf("Found %d matches:\n", len(matches))
			for _, match := range matches {
				fmt.Printf("  %s\n", match)
			}
		} else {
			fmt.Println("No matches found")
		}
	}

	// Demonstrate prefix search with manual regex filtering
	fmt.Println("\n=== Prefix Search + Regex Filter ===")
	prefix := "app"
	pattern := ".*ion$"
	
	fmt.Printf("Prefix: %s, Pattern: %s\n", prefix, pattern)
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Error compiling regex: %v", err)
	}

	matches := []string{}
	prefixLower := strings.ToLower(prefix)
	
	// Use FST prefix iterator for efficiency
	iter := fstIndex.PrefixIterator([]byte(prefixLower))
	for iter.HasNext() {
		key, docID := iter.Next()
		keyStr := string(key)
		
		if regex.MatchString(keyStr) {
			if int(docID) < len(documents) && docID > 0 {
				docContent := documents[docID-1]
				matches = append(matches, docContent)
			}
		}
	}

	if len(matches) > 0 {
		fmt.Printf("Found %d matches:\n", len(matches))
		for _, match := range matches {
			fmt.Printf("  %s\n", match)
		}
	} else {
		fmt.Println("No matches found")
	}

	// Show FST capabilities
	fmt.Println("\n=== FST Basic Capabilities ===")
	fmt.Println("\n1. Exact lookup:")
	if docID, found := fstIndex.Get([]byte("application")); found {
		fmt.Printf("  'application' -> Document ID: %d\n", docID)
		if int(docID) <= len(documents) && docID > 0 {
			fmt.Printf("  Document: %s\n", documents[docID-1])
		}
	}

	fmt.Println("\n2. Prefix search:")
	iter = fstIndex.PrefixIterator([]byte("app"))
	fmt.Println("  Words starting with 'app':")
	for iter.HasNext() {
		key, docID := iter.Next()
		if int(docID) <= len(documents) && docID > 0 {
			fmt.Printf("    %s -> %s\n", string(key), documents[docID-1])
		}
	}

	fmt.Println("\n3. Regular search engine:")
	results := searchEngine.Search("app")
	fmt.Printf("  Search 'app' found %d results:\n", len(results))
	for _, result := range results {
		fmt.Printf("    %s (score: %.2f)\n", result.Word, result.Score)
	}

	fmt.Printf("\n=== Demo Complete ===\n")
	fmt.Println("Note: Full regex integration is coming soon!")
	fmt.Println("This demo shows how FST + Go regex can work together.")
}
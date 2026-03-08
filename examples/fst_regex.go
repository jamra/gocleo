package main

import (
	"fmt"
	"log"

	"github.com/jamra/gocleo/internal/fst"
	"github.com/jamra/gocleo/internal/scoring"
)

func main() {
	fmt.Println("=== FST Regex Search Example ===")

	// Sample documents with various patterns
	documents := []string{
		"apple pie recipe",
		"banana bread instructions", 
		"cherry cake tutorial",
		"date pudding guide",
		"elderberry jam making",
		"fig tart preparation",
		"grape wine brewing",
		"honey cake baking",
		"ice cream making",
		"jelly donut recipe",
		"kiwi fruit salad",
		"lemon tart recipe",
		"mango smoothie blend",
		"nutmeg cookies baking",
		"orange juice fresh",
		"peach cobbler dessert",
		"quince preserve making",
		"raspberry pie filling",
		"strawberry shortcake",
		"tomato sauce cooking",
		"vanilla extract making",
		"watermelon juice blend",
		"application development",
		"programming tutorial",
		"development environment",
		"testing framework",
		"debugging techniques",
		"performance optimization",
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

	// Test basic regex patterns
	testRegexPatterns := []string{
		"app.*",           // Words starting with "app"
		".*ing$",          // Words ending with "ing"
		"^[a-c].*",        // Words starting with a, b, or c
		".*[aeiou]{2}.*",  // Words containing double vowels
		"^.{4}$",          // Exactly 4 characters
		".*recipe.*",      // Words containing "recipe"
		"^(cake|pie)$",    // Exactly "cake" or "pie"
		".*[0-9].*",       // Words containing numbers
		"^[^aeiou].*",     // Words starting with consonants
		".*ment$",         // Words ending with "ment"
	}

	fmt.Println("=== Basic Regex Search Tests ===")
	for _, pattern := range testRegexPatterns {
		fmt.Printf("\nRegex Pattern: %s\n", pattern)
		results, err := searchEngine.RegexSearch(pattern)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if len(results) > 0 {
			fmt.Printf("Found %d matches:\n", len(results))
			for i, result := range results {
				if i >= 5 { // Limit to first 5 results
					fmt.Printf("  ... and %d more\n", len(results)-i)
					break
				}
				fmt.Printf("  %s (score: %.2f)\n", result.Word, result.Score)
			}
		} else {
			fmt.Println("No matches found")
		}
	}

	// Test prefix + regex combination
	fmt.Println("\n=== Prefix + Regex Search Tests ===")
	prefixRegexTests := []struct {
		prefix  string
		pattern string
		desc    string
	}{
		{"app", ".*e$", "Words with prefix 'app' ending with 'e'"},
		{"recipe", ".*", "All words with prefix 'recipe'"},
		{"baking", ".*", "All words with prefix 'baking'"},
		{"dev", ".*ment$", "Words with prefix 'dev' ending with 'ment'"},
		{"test", ".*ing$", "Words with prefix 'test' ending with 'ing'"},
	}

	for _, test := range prefixRegexTests {
		fmt.Printf("\n%s\n", test.desc)
		fmt.Printf("Prefix: %s, Regex: %s\n", test.prefix, test.pattern)
		
		results, err := searchEngine.PrefixRegexSearch(test.prefix, test.pattern)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if len(results) > 0 {
			fmt.Printf("Found %d matches:\n", len(results))
			for _, result := range results {
				fmt.Printf("  %s (score: %.2f)\n", result.Word, result.Score)
			}
		} else {
			fmt.Println("No matches found")
		}
	}

	// Test complex search with multiple criteria
	fmt.Println("\n=== Complex Search Tests ===")
	complexTests := []fst.ComplexSearchOptions{
		{
			Prefix:       "app",
			RegexPattern: ".*ion$",
			Limit:        5,
		},
		{
			RegexPattern:     ".*ing$",
			FuzzyPattern:     "baking",
			FuzzyMaxDistance: 2,
			Limit:            3,
		},
		{
			Prefix:       "rec",
			RegexPattern: ".*",
			Limit:        10,
		},
	}

	for i, options := range complexTests {
		fmt.Printf("\nComplex Search %d:\n", i+1)
		fmt.Printf("Options: Prefix='%s', Regex='%s', Fuzzy='%s' (dist=%d), Limit=%d\n", 
			options.Prefix, options.RegexPattern, options.FuzzyPattern, options.FuzzyMaxDistance, options.Limit)
		
		results, err := searchEngine.ComplexSearch(options)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if len(results) > 0 {
			fmt.Printf("Found %d matches:\n", len(results))
			for _, result := range results {
				fmt.Printf("  %s (score: %.2f)\n", result.Word, result.Score)
			}
		} else {
			fmt.Println("No matches found")
		}
	}

	// Test some common use cases
	fmt.Println("\n=== Common Use Cases ===")
	
	fmt.Println("\n1. Find all food items (words ending with common food terms):")
	foodPattern := ".*(cake|pie|bread|juice|sauce|jam|tart).*"
	results, _ := searchEngine.RegexSearch(foodPattern)
	for _, result := range results[:min(5, len(results))] {
		fmt.Printf("  %s\n", result.Word)
	}

	fmt.Println("\n2. Find all action words (ending with -ing):")
	actionPattern := ".*ing$"
	results, _ = searchEngine.RegexSearch(actionPattern)
	for _, result := range results[:min(5, len(results))] {
		fmt.Printf("  %s\n", result.Word)
	}

	fmt.Println("\n3. Find all technical terms (containing 'dev', 'prog', or 'test'):")
	techPattern := ".*(dev|prog|test).*"
	results, _ = searchEngine.RegexSearch(techPattern)
	for _, result := range results[:min(5, len(results))] {
		fmt.Printf("  %s\n", result.Word)
	}

	fmt.Printf("\n=== FST Regex Search Demo Complete ===\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
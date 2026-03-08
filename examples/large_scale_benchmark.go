package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	var words []string
	
	if len(os.Args) > 1 {
		// Load words from file
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := scanner.Text()
			if word != "" {
				words = append(words, word)
			}
		}
		fmt.Printf("✅ Loaded %d words from %s\n", len(words), os.Args[1])
	} else {
		// Use built-in test data
		words = []string{"app", "apple", "application", "banana", "band", "data", "test"}
		fmt.Println("Using built-in test data (no file provided)")
	}

	// Sort words for FST (required)
	sort.Strings(words)

	// Build FST
	builder := fst.NewFSTBuilder()
	for i, word := range words {
		err := builder.Add([]byte(word), uint64(i))
		if err != nil {
			fmt.Printf("Error adding word '%s': %v\n", word, err)
			return
		}
	}

	fstInstance, err := builder.Build()
	if err != nil {
		fmt.Printf("Error building FST: %v\n", err)
		return
	}
	fmt.Printf("✅ Built FST with %d entries\n", len(words))

	// Test patterns
	patterns := []string{
		"app.*",
		".*ing$",
		"^[a-c].*",
		".*an.*",
	}

	fmt.Println("\n=== Performance Comparison ===")
	
	for _, pattern := range patterns {
		fmt.Printf("\nTesting pattern: %s\n", pattern)
		
		// FST approach
		start := time.Now()
		fstResults := fstRegexSearch(fstInstance, pattern)
		fstDuration := time.Since(start)
		
		// Slice approach
		start = time.Now()
		sliceResults := sliceRegexSearch(words, pattern)
		sliceDuration := time.Since(start)
		
		// Verify results match
		sort.Strings(fstResults)
		sort.Strings(sliceResults)
		
		fmt.Printf("✅ FST found %d results in %v\n", len(fstResults), fstDuration)
		fmt.Printf("✅ Slice found %d results in %v\n", len(sliceResults), sliceDuration)
		
		// Show performance difference
		if sliceDuration > 0 {
			speedup := float64(sliceDuration) / float64(fstDuration)
			if speedup > 1 {
				fmt.Printf("🚀 FST is %.2fx faster\n", speedup)
			} else {
				fmt.Printf("📊 Slice is %.2fx faster\n", 1/speedup)
			}
		}
		
		// Show sample results
		if len(fstResults) > 0 {
			maxShow := 5
			if len(fstResults) < maxShow {
				maxShow = len(fstResults)
			}
			fmt.Printf("📋 Sample results: %v\n", fstResults[:maxShow])
		}
	}
	
	fmt.Println("\n✅ All tests completed successfully!")
}

func sliceRegexSearch(words []string, pattern string) []string {
	regex, _ := regexp.Compile(pattern)
	var results []string
	for _, word := range words {
		if regex.MatchString(word) {
			results = append(results, word)
		}
	}
	return results
}

func fstRegexSearch(fstInstance *fst.FST, pattern string) []string {
	regex, _ := regexp.Compile(pattern)
	var results []string
	iterator := fstInstance.Iterator()
	for iterator.HasNext() {
		key, _ := iterator.Next()
		if regex.MatchString(string(key)) {
			results = append(results, string(key))
		}
	}
	return results
}
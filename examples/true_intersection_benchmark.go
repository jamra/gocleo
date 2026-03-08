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
		words = []string{"app", "apple", "application", "banana", "band", "data", "test", "testing", "tester"}
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
	}

	fmt.Println("\n=== TRUE AUTOMATA INTERSECTION vs NAIVE APPROACHES ===")
	
	for _, pattern := range patterns {
		fmt.Printf("\n🧪 Testing pattern: %s\n", pattern)
		
		// Method 1: Naive slice + regex (baseline)
		start := time.Now()
		sliceResults := sliceRegexSearch(words, pattern)
		sliceDuration := time.Since(start)
		
		// Method 2: FST iteration + regex (our previous approach)
		start = time.Now()
		fstResults := fstRegexSearch(fstInstance, pattern)
		fstDuration := time.Since(start)
		
		// Method 3: TRUE automata intersection (the holy grail!)
		start = time.Now()
		trueResults, err := trueAutomataIntersection(fstInstance, pattern)
		trueDuration := time.Since(start)
		if err != nil {
			fmt.Printf("❌ True intersection failed: %v\n", err)
			continue
		}
		
		// Verify all methods return same results
		sort.Strings(sliceResults)
		sort.Strings(fstResults)
		sort.Strings(trueResults)
		
		fmt.Printf("📊 Results comparison:\n")
		fmt.Printf("  Slice:           %d results in %v\n", len(sliceResults), sliceDuration)
		fmt.Printf("  FST iteration:   %d results in %v\n", len(fstResults), fstDuration)
		fmt.Printf("  True intersection: %d results in %v\n", len(trueResults), trueDuration)
		
		// Show performance ratios
		if sliceDuration > 0 {
			fstRatio := float64(fstDuration) / float64(sliceDuration)
			trueRatio := float64(trueDuration) / float64(sliceDuration)
			
			fmt.Printf("📈 Performance vs slice baseline:\n")
			if fstRatio > 1 {
				fmt.Printf("  FST iteration:   %.2fx SLOWER\n", fstRatio)
			} else {
				fmt.Printf("  FST iteration:   %.2fx faster\n", 1/fstRatio)
			}
			
			if trueRatio > 1 {
				fmt.Printf("  True intersection: %.2fx SLOWER\n", trueRatio)
			} else {
				fmt.Printf("  True intersection: %.2fx FASTER 🚀\n", 1/trueRatio)
			}
		}
		
		// Show sample results
		if len(trueResults) > 0 {
			sampleSize := 5
			if len(trueResults) < sampleSize {
				sampleSize = len(trueResults)
			}
			fmt.Printf("📋 Sample results: %v\n", trueResults[:sampleSize])
		}
	}
	
	fmt.Println("\n✅ TRUE automata intersection benchmark completed!")
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
		keyStr := string(key)
		if regex.MatchString(keyStr) {
			results = append(results, keyStr)
		}
	}
	return results
}

func trueAutomataIntersection(fstInstance *fst.FST, pattern string) ([]string, error) {
	// Create true regex automaton
	regexAutomaton, err := fst.NewTrueRegexAutomaton(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create regex automaton: %v", err)
	}
	
	// Perform TRUE automata intersection (not iteration!)
	return regexAutomaton.TrueAutomataIntersection(fstInstance)
}
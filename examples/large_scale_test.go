package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run large_scale_test.go <word_list_file>")
		fmt.Println("Example: go run large_scale_test.go words_479k.txt")
		os.Exit(1)
	}

	wordFile := os.Args[1]
	fmt.Printf("=== Large Scale FST Test ===\n")
	fmt.Printf("Loading words from: %s\n", wordFile)

	// Load words from file
	words, err := loadWordsFromFile(wordFile)
	if err != nil {
		fmt.Printf("Error loading words: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Loaded %d words\n", len(words))
	sort.Strings(words)

	// Build FST
	fmt.Printf("Building FST...\n")
	start := time.Now()
	builder := fst.NewFSTBuilder()
	for i, word := range words {
		builder.Add([]byte(word), uint64(i))
	}
	fstInstance, err := builder.Build()
	if err != nil {
		fmt.Printf("Error building FST: %v\n", err)
		os.Exit(1)
	}
	buildTime := time.Since(start)
	fmt.Printf("✅ FST built in %v\n", buildTime)

	// Test patterns
	patterns := []string{"app.*", ".*ing$", "^[a-c].*"}

	for _, pattern := range patterns {
		fmt.Printf("\nTesting pattern: %s\n", pattern)
		
		// Slice approach
		start = time.Now()
		sliceResults := sliceRegexSearch(words, pattern)
		sliceTime := time.Since(start)
		
		// FST approach
		start = time.Now()
		fstResults := fstRegexSearch(fstInstance, pattern)
		fstTime := time.Since(start)
		
		fmt.Printf("  Slice: %d results in %v\n", len(sliceResults), sliceTime)
		fmt.Printf("  FST: %d results in %v\n", len(fstResults), fstTime)
		
		if sliceTime > fstTime {
			speedup := float64(sliceTime) / float64(fstTime)
			fmt.Printf("  → FST is %.2fx faster! 🚀\n", speedup)
		} else {
			slowdown := float64(fstTime) / float64(sliceTime)
			fmt.Printf("  → Slice is %.2fx faster\n", slowdown)
		}
	}
}

func loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	return words, scanner.Err()
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
	for iterator.Next() {
		key := string(iterator.Current())
		if regex.MatchString(key) {
			results = append(results, key)
		}
	}
	return results
}
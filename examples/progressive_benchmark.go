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
	if len(os.Args) < 2 {
		// Test with smaller built-in data first
		testWithBuiltInData()
		return
	}

	filename := os.Args[1]
	fmt.Printf("=== Progressive FST Benchmark ===\n")
	fmt.Printf("📁 Loading words from: %s\n", filename)

	// Load words
	words, err := loadWords(filename)
	if err != nil {
		fmt.Printf("❌ Error loading words: %v\n", err)
		return
	}

	fmt.Printf("✅ Loaded %d words\n", len(words))

	// Test with progressively larger datasets
	testSizes := []int{1000, 5000, 10000, 25000, 50000}
	
	// Add the full size if it's not too huge
	if len(words) <= 100000 {
		testSizes = append(testSizes, len(words))
	} else {
		testSizes = append(testSizes, 100000)
		fmt.Printf("⚠️  Dataset is huge (%d words). Testing up to 100k to avoid timeout.\n", len(words))
	}

	for _, size := range testSizes {
		if size > len(words) {
			continue
		}
		fmt.Printf("\n🧪 Testing with %d words...\n", size)
		testSubset := words[:size]
		
		// Run the benchmark
		runBenchmark(testSubset, size)
	}
}

func runBenchmark(words []string, size int) {
	patterns := []string{
		"app.*",     // Simple prefix
		".*ing$",    // Simple suffix  
		"^[a-c].*",  // Character class prefix
	}

	fmt.Printf("📊 Building FST with %d words...", size)
	start := time.Now()
	
	// Build FST
	builder := fst.NewFSTBuilder()
	sort.Strings(words) // FST requires sorted input
	
	for i, word := range words {
		builder.Insert([]byte(word), uint64(i))
	}
	
	testFST, err := builder.Build()
	if err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		return
	}
	
	buildTime := time.Since(start)
	fmt.Printf(" ✅ Built in %v\n", buildTime)

	// Test each pattern
	for _, pattern := range patterns {
		fmt.Printf("  🔍 Testing pattern: %s\n", pattern)
		
		// FST approach
		fstStart := time.Now()
		fstResults := searchFST(testFST, pattern)
		fstTime := time.Since(fstStart)
		
		// Slice approach
		sliceStart := time.Now()
		sliceResults := searchSlice(words, pattern)
		sliceTime := time.Since(sliceStart)
		
		// Compare
		if len(fstResults) == len(sliceResults) {
			fmt.Printf("    ✅ Both found %d results\n", len(fstResults))
		} else {
			fmt.Printf("    ⚠️  Results differ! FST:%d, Slice:%d\n", len(fstResults), len(sliceResults))
		}
		
		if fstTime < sliceTime {
			speedup := float64(sliceTime) / float64(fstTime)
			fmt.Printf("    🚀 FST wins! %v vs %v (%.2fx faster)\n", fstTime, sliceTime, speedup)
		} else {
			slowdown := float64(fstTime) / float64(sliceTime)
			fmt.Printf("    📈 Slice wins: %v vs %v (%.2fx slower)\n", fstTime, sliceTime, slowdown)
		}
	}
}

func searchFST(testFST *fst.FST, pattern string) []string {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	var results []string
	iterator := testFST.Iterator()
	
	for iterator.HasNext() {
		key, _ := iterator.Next()
		keyStr := string(key)
		if regex.MatchString(keyStr) {
			results = append(results, keyStr)
		}
	}
	
	return results
}

func searchSlice(words []string, pattern string) []string {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	var results []string
	for _, word := range words {
		if regex.MatchString(word) {
			results = append(results, word)
		}
	}
	
	return results
}

func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			words = append(words, word)
		}
	}
	
	return words, scanner.Err()
}

func testWithBuiltInData() {
	fmt.Println("=== Testing with built-in data ===")
	words := []string{"app", "apple", "application", "banana", "band", "test", "testing"}
	runBenchmark(words, len(words))
}
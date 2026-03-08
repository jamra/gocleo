package fst

import (
	"fmt"
	"regexp"
	"sort"
	"testing"
	"time"
)

// TestComprehensiveRegexAccuracy tests accuracy of FST + regex approach
func TestComprehensiveRegexAccuracy(t *testing.T) {
	// Create comprehensive test dataset
	testWords := generateComprehensiveTestSet()
	
	// Build FST
	builder := NewFSTBuilder()
	sort.Strings(testWords) // Ensure lexicographic order
	
	for i, word := range testWords {
		err := builder.Add([]byte(word), uint64(i+1))
		if err != nil {
			t.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fst, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FST: %v", err)
	}

	// Test patterns with increasing complexity
	testCases := []struct {
		name     string
		pattern  string
		testType string
	}{
		{"Simple literal", "test", "exact"},
		{"Simple suffix", ".*ing$", "suffix"},
		{"Simple prefix", "^app.*", "prefix"}, 
		{"Contains", ".*data.*", "contains"},
		{"Character class", "[a-c].*", "char_class"},
		{"Alternation", "(test|exam).*", "alternation"},
		{"Quantifier", "a+b*", "quantifier"},
		{"Complex", "^(data|info).*base.*$", "complex"},
		{"Anchored both", "^test.*ing$", "anchored"},
		{"Optional", "colou?r", "optional"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get FST results
			fstResults, fstTime := getFSTRegexResults(fst, tc.pattern)
			
			// Get expected results (baseline truth)
			expectedResults, baselineTime := getExpectedResults(testWords, tc.pattern)
			
			// Compare accuracy
			if !compareStringSlices(fstResults, expectedResults) {
				t.Errorf("Accuracy mismatch for pattern '%s'", tc.pattern)
				t.Errorf("FST found: %d results", len(fstResults))
				t.Errorf("Expected: %d results", len(expectedResults))
				
				// Show differences
				fstMap := stringSliceToMap(fstResults)
				expectedMap := stringSliceToMap(expectedResults)
				
				for result := range expectedMap {
					if !fstMap[result] {
						t.Errorf("Missing: %s", result)
					}
				}
				
				for result := range fstMap {
					if !expectedMap[result] {
						t.Errorf("Extra: %s", result)
					}
				}
			} else {
				t.Logf("✅ Pattern '%s' (%s): %d matches, FST: %v, Baseline: %v", 
					tc.pattern, tc.testType, len(fstResults), fstTime, baselineTime)
			}
		})
	}
}

// BenchmarkRegexComplexity benchmarks different regex complexity levels
func BenchmarkRegexComplexity(b *testing.B) {
	// Create test FST
	testWords := generateComprehensiveTestSet()
	builder := NewFSTBuilder()
	sort.Strings(testWords)
	
	for i, word := range testWords {
		err := builder.Add([]byte(word), uint64(i+1))
		if err != nil {
			b.Fatalf("Failed to add word: %v", err)
		}
	}

	fst, err := builder.Build()
	if err != nil {
		b.Fatalf("Failed to build FST: %v", err)
	}

	complexityLevels := map[string]string{
		"Simple":     "test",
		"Suffix":     ".*ing$", 
		"Prefix":     "^app.*",
		"Contains":   ".*data.*",
		"CharClass":  "[a-z]+er$",
		"Alternation": "(test|data|app).*",
		"Quantifier": "a+b*c?",
		"Complex":    "^(data|info).*base.*$",
	}

	for name, pattern := range complexityLevels {
		b.Run(name, func(b *testing.B) {
			regex, err := regexp.Compile(pattern)
			if err != nil {
				b.Skipf("Invalid pattern: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var matches []string
				iterator := fst.Iterator()
				for iterator.HasNext() {
					key, _ := iterator.Next()
					keyStr := string(key)
					if regex.MatchString(keyStr) {
						matches = append(matches, keyStr)
					}
				}
				_ = matches
			}
		})
	}
}

// BenchmarkDatasetSize benchmarks performance scaling with dataset size
func BenchmarkDatasetSize(b *testing.B) {
	sizes := []int{100, 500, 1000, 5000}
	pattern := ".*ing$"
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			// Create FST of specific size
			words := generateTestWords(size)
			builder := NewFSTBuilder()
			sort.Strings(words)
			
			for i, word := range words {
				builder.Add([]byte(word), uint64(i+1))
			}

			fst, err := builder.Build()
			if err != nil {
				b.Fatalf("Failed to build FST: %v", err)
			}

			regex, err := regexp.Compile(pattern)
			if err != nil {
				b.Fatalf("Failed to compile regex: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var matches []string
				iterator := fst.Iterator()
				for iterator.HasNext() {
					key, _ := iterator.Next()
					keyStr := string(key)
					if regex.MatchString(keyStr) {
						matches = append(matches, keyStr)
					}
				}
			}
			
			b.ReportMetric(float64(size), "dataset_size")
		})
	}
}

// BenchmarkApproachComparison compares different search approaches
func BenchmarkApproachComparison(b *testing.B) {
	// Create test data
	testWords := generateTestWords(1000)
	sort.Strings(testWords)
	
	// Build FST
	builder := NewFSTBuilder()
	for i, word := range testWords {
		builder.Add([]byte(word), uint64(i+1))
	}

	fst, err := builder.Build()
	if err != nil {
		b.Fatalf("Failed to build FST: %v", err)
	}

	pattern := ".*ing$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		b.Fatalf("Failed to compile regex: %v", err)
	}

	// Approach 1: FST iteration + regex
	b.Run("FST_Plus_Regex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var matches []string
			iterator := fst.Iterator()
			for iterator.HasNext() {
				key, _ := iterator.Next()
				keyStr := string(key)
				if regex.MatchString(keyStr) {
					matches = append(matches, keyStr)
				}
			}
		}
	})

	// Approach 2: Direct slice + regex (baseline)
	b.Run("Slice_Plus_Regex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var matches []string
			for _, word := range testWords {
				if regex.MatchString(word) {
					matches = append(matches, word)
				}
			}
		}
	})
	
	// Approach 3: FST prefix + regex (when applicable)
	prefixPattern := "^app.*ing$"
	prefixRegex, _ := regexp.Compile(prefixPattern)
	
	b.Run("FST_Prefix_Plus_Regex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var matches []string
			iterator := fst.PrefixIterator([]byte("app"))
			for iterator.HasNext() {
				key, _ := iterator.Next()
				keyStr := string(key)
				if prefixRegex.MatchString(keyStr) {
					matches = append(matches, keyStr)
				}
			}
		}
	})
}

// Helper functions

func generateComprehensiveTestSet() []string {
	words := []string{}
	
	// Common English words
	common := []string{"test", "data", "app", "run", "walk", "code", "file", "user", "search", "build"}
	suffixes := []string{"", "ing", "ed", "er", "s", "able", "tion", "ment", "ly", "ness"}
	prefixes := []string{"", "pre", "un", "re", "dis", "over", "under", "out"}
	
	for _, prefix := range prefixes {
		for _, root := range common {
			for _, suffix := range suffixes {
				word := prefix + root + suffix
				if len(word) > 0 {
					words = append(words, word)
				}
			}
		}
	}
	
	// Add specific test cases
	specific := []string{
		"apple", "application", "apply", "appreciate",
		"database", "datafile", "metadata", "information", 
		"running", "walking", "jumping", "testing", "coding",
		"color", "colour", "honor", "honour",
		"exam", "example", "execute", "exercise",
	}
	
	words = append(words, specific...)
	
	// Remove duplicates
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		if len(word) > 0 {
			uniqueWords[word] = true
		}
	}
	
	result := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		result = append(result, word)
	}
	
	return result
}

func generateTestWords(count int) []string {
	baseWords := []string{"test", "data", "app", "run", "walk", "code", "search", "build", "file", "user", "info", "create", "update", "delete"}
	suffixes := []string{"", "ing", "ed", "er", "s", "able", "tion", "ment", "ly", "ness"}
	
	words := []string{}
	for i := 0; i < count && len(words) < count; i++ {
		base := baseWords[i%len(baseWords)]
		suffix := suffixes[(i/len(baseWords))%len(suffixes)]
		word := fmt.Sprintf("%s%s%d", base, suffix, i/100) // Add number to ensure uniqueness
		words = append(words, word)
	}
	
	return words[:count]
}

func getFSTRegexResults(fst *FST, pattern string) ([]string, time.Duration) {
	start := time.Now()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, time.Since(start)
	}

	var results []string
	iterator := fst.Iterator()
	for iterator.HasNext() {
		key, _ := iterator.Next()
		keyStr := string(key)
		if regex.MatchString(keyStr) {
			results = append(results, keyStr)
		}
	}
	
	sort.Strings(results)
	return results, time.Since(start)
}

func getExpectedResults(words []string, pattern string) ([]string, time.Duration) {
	start := time.Now()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, time.Since(start)
	}

	var results []string
	for _, word := range words {
		if regex.MatchString(word) {
			results = append(results, word)
		}
	}
	
	sort.Strings(results)
	return results, time.Since(start)
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	sortedA := make([]string, len(a))
	sortedB := make([]string, len(b))
	copy(sortedA, a)
	copy(sortedB, b)
	
	sort.Strings(sortedA)
	sort.Strings(sortedB)
	
	for i := range sortedA {
		if sortedA[i] != sortedB[i] {
			return false
		}
	}
	
	return true
}

func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range slice {
		result[s] = true
	}
	return result
}
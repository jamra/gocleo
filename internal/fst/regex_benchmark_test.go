package fst

import (
	"regexp"
	"sort"
	"testing"
)

// TestRegexWithFST verifies regex functionality works with FST
func TestRegexWithFST(t *testing.T) {
	builder := NewFSTBuilder()
	
	testData := []string{
		"running", "walking", "jumping", 
		"test", "testing", "tester",
		"app", "apple", "application",
		"data", "database", "datafile",
	}
	
	// Sort the test data for lexicographic order requirement
	sort.Strings(testData)
	
	for i, word := range testData {
		err := builder.Add([]byte(word), uint64(i+1))
		if err != nil {
			t.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fst, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FST: %v", err)
	}

	// Test pattern matching
	pattern := ".*ing$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
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

	if len(results) == 0 {
		t.Fatal("No matches found - regex integration failed")
	}

	t.Logf("✅ FST + Regex integration working! Found %d matches for pattern '%s'", len(results), pattern)
	t.Logf("Matches: %v", results)
}

// BenchmarkFSTRegexSearch benchmarks FST iteration + regex matching
func BenchmarkFSTRegexSearch(b *testing.B) {
	// Create FST with test data
	builder := NewFSTBuilder()
	
	// Generate realistic dataset
	prefixes := []string{"app", "data", "test", "run", "walk", "code", "file", "user", "search", "build"}
	suffixes := []string{"", "ing", "ed", "er", "s", "able", "tion", "ment", "ly", "ness"}
	
	var allWords []string
	for _, prefix := range prefixes {
		for _, suffix := range suffixes {
			word := prefix + suffix
			allWords = append(allWords, word)
		}
	}
	
	// Sort for lexicographic order
	sort.Strings(allWords)
	
	for i, word := range allWords {
		err := builder.Add([]byte(word), uint64(i+1))
		if err != nil {
			b.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fst, err := builder.Build()
	if err != nil {
		b.Fatalf("Failed to build FST: %v", err)
	}

	// Test different patterns
	patterns := map[string]*regexp.Regexp{
		"suffix_ing":     regexp.MustCompile(".*ing$"),
		"prefix_app":     regexp.MustCompile("^app.*"),
		"contains_data":  regexp.MustCompile(".*data.*"),
		"suffix_ed":      regexp.MustCompile(".*ed$"),
		"prefix_test":    regexp.MustCompile("^test.*"),
	}

	for name, regex := range patterns {
		b.Run(name, func(b *testing.B) {
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
				_ = matches // Prevent optimization
			}
		})
	}
}

// BenchmarkRegexCompilation benchmarks regex compilation overhead
func BenchmarkRegexCompilation(b *testing.B) {
	patterns := []string{
		".*ing$",
		"^app.*",
		".*data.*",
		"[a-z]+ed$",
		"^(test|exam|quiz).*",
	}

	for _, pattern := range patterns {
		b.Run(pattern, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := regexp.Compile(pattern)
				if err != nil {
					b.Fatalf("Failed to compile pattern %s: %v", pattern, err)
				}
			}
		})
	}
}

// BenchmarkMemoryUsage benchmarks memory usage during regex search
func BenchmarkMemoryUsage(b *testing.B) {
	// Create FST
	builder := NewFSTBuilder()
	var allWords []string
	for i := 0; i < 1000; i++ {
		word := generateTestWord(i)
		allWords = append(allWords, word)
	}
	
	// Remove duplicates and sort
	uniqueWords := make(map[string]bool)
	for _, word := range allWords {
		uniqueWords[word] = true
	}
	
	var sortedWords []string
	for word := range uniqueWords {
		sortedWords = append(sortedWords, word)
	}
	sort.Strings(sortedWords)
	
	for i, word := range sortedWords {
		err := builder.Add([]byte(word), uint64(i+1))
		if err != nil {
			b.Fatalf("Failed to add word '%s': %v", word, err)
		}
	}

	fst, err := builder.Build()
	if err != nil {
		b.Fatalf("Failed to build FST: %v", err)
	}

	regex := regexp.MustCompile(".*ing$")

	b.Run("MemoryAllocs", func(b *testing.B) {
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
}

// Helper function to generate test words
func generateTestWord(seed int) string {
	bases := []string{"test", "data", "app", "run", "walk", "code", "build", "search", "file", "user"}
	suffixes := []string{"", "ing", "ed", "er", "s", "able", "tion", "ly"}
	
	base := bases[seed%len(bases)]
	suffix := suffixes[(seed/len(bases))%len(suffixes)]
	
	return base + suffix
}
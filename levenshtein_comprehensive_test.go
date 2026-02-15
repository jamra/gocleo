package cleo

import "testing"

func TestLevenshteinDistanceComprehensive(t *testing.T) {
    tests := []struct {
        s1, s2   string
        expected int
        name     string
    }{
        // Basic cases
        {"", "", 0, "empty strings"},
        {"", "a", 1, "empty to single char"},
        {"a", "", 1, "single char to empty"},
        {"a", "a", 0, "identical single char"},
        
        // Single character operations
        {"cat", "bat", 1, "single substitution"},
        {"cat", "cats", 1, "single insertion"},
        {"cats", "cat", 1, "single deletion"},
        
        // Classic examples
        {"kitten", "sitting", 3, "classic kitten->sitting"},
        {"saturday", "sunday", 3, "saturday->sunday"},
        
        // Multiple operations
        {"abc", "def", 3, "all substitutions"},
        {"hello", "world", 4, "multiple operations"},
        {"algorithm", "logarithm", 3, "prefix difference"},
        
        // Longer strings
        {"programming", "algorithm", 10, "longer different strings"},
        {"levenshtein", "distance", 10, "algorithm name vs concept"},
        {"javascript", "typescript", 4, "similar languages"},
        
        // Pattern cases
        {"abab", "baba", 2, "pattern shift"},
        {"aaa", "aaaa", 1, "repeated chars insertion"},
        {"aaaa", "aaa", 1, "repeated chars deletion"},
        
        // Edge cases
        {"a", "ab", 1, "single to double"},
        {"ab", "a", 1, "double to single"},
        {"abc", "ac", 1, "middle deletion"},
        {"ac", "abc", 1, "middle insertion"},
        {"abcd", "efgh", 4, "complete replacement"},
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := LevenshteinDistance(test.s1, test.s2)
            if result != test.expected {
                t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", 
                    test.s1, test.s2, result, test.expected)
            }
        })
    }
}

// Property-based tests to ensure mathematical properties hold
func TestLevenshteinProperties(t *testing.T) {
    t.Run("symmetry", func(t *testing.T) {
        // d(s1, s2) should equal d(s2, s1)
        testCases := [][2]string{
            {"hello", "world"},
            {"cat", "dog"},
            {"algorithm", "logarithm"},
            {"", "test"},
            {"abc", "xyz"},
        }
        
        for _, tc := range testCases {
            s1, s2 := tc[0], tc[1]
            d1 := LevenshteinDistance(s1, s2)
            d2 := LevenshteinDistance(s2, s1)
            if d1 != d2 {
                t.Errorf("Symmetry violated: d(%q, %q) = %d, d(%q, %q) = %d", 
                    s1, s2, d1, s2, s1, d2)
            }
        }
    })

    t.Run("identity", func(t *testing.T) {
        // d(s, s) should always be 0
        strings := []string{"", "a", "hello", "longer string", "test123"}
        for _, s := range strings {
            result := LevenshteinDistance(s, s)
            if result != 0 {
                t.Errorf("Identity violated: d(%q, %q) = %d, want 0", s, s, result)
            }
        }
    })

    t.Run("triangle inequality", func(t *testing.T) {
        // d(s1, s3) <= d(s1, s2) + d(s2, s3)
        testCases := [][3]string{
            {"cat", "bat", "rat"},
            {"hello", "help", "held"},
            {"abc", "ab", "a"},
            {"test", "best", "rest"},
        }
        
        for _, tc := range testCases {
            s1, s2, s3 := tc[0], tc[1], tc[2]
            d13 := LevenshteinDistance(s1, s3)
            d12 := LevenshteinDistance(s1, s2)
            d23 := LevenshteinDistance(s2, s3)
            
            if d13 > d12 + d23 {
                t.Errorf("Triangle inequality violated: d(%q,%q)=%d > d(%q,%q)=%d + d(%q,%q)=%d", 
                    s1, s3, d13, s1, s2, d12, s2, s3, d23)
            }
        }
    })

    t.Run("non-negativity", func(t *testing.T) {
        // Distance should never be negative
        testCases := [][2]string{
            {"", ""},
            {"a", "b"},
            {"hello", "world"},
            {"test", ""},
            {"", "test"},
        }
        
        for _, tc := range testCases {
            s1, s2 := tc[0], tc[1]
            result := LevenshteinDistance(s1, s2)
            if result < 0 {
                t.Errorf("Non-negativity violated: d(%q, %q) = %d", s1, s2, result)
            }
        }
    })
}

// Benchmark tests for performance analysis
func BenchmarkLevenshteinDistance(b *testing.B) {
    benchmarks := []struct {
        name string
        s1, s2 string
    }{
        {"empty", "", ""},
        {"short", "cat", "bat"},
        {"medium", "hello world", "jello world"},
        {"long", "this is a longer string for testing", "this is a different longer string for testing"},
        {"very_long", 
            "this is a very long string that we will use to test the performance of our levenshtein distance implementation with many characters", 
            "this is a very long string that we will use to test the performance of our levenstein distance implementation with many characters"},
        {"different_lengths", "short", "much longer string here"},
        {"empty_vs_long", "", "this is a reasonably long string"},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                _ = LevenshteinDistance(bm.s1, bm.s2)
            }
        })
    }
}

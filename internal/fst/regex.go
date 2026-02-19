package fst

import (
	"regexp"
)

// RegexMatcher provides regex matching capabilities for FSAs
type RegexMatcher struct {
	pattern *regexp.Regexp
}

// NewRegexMatcher creates a new regex matcher
func NewRegexMatcher(pattern string) (*RegexMatcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	
	return &RegexMatcher{
		pattern: re,
	}, nil
}

// Match tests if the regex matches the given input
func (rm *RegexMatcher) Match(input []byte) bool {
	return rm.pattern.Match(input)
}

// FindMatches finds all matches in the given strings
func (rm *RegexMatcher) FindMatches(keys []string) []string {
	var matches []string
	
	for _, key := range keys {
		if rm.pattern.MatchString(key) {
			matches = append(matches, key)
		}
	}
	
	return matches
}

// RegexAutomaton represents a simple regex-based automaton
// This is a simplified implementation - a full implementation would
// convert regex to NFA/DFA for better performance
type RegexAutomaton struct {
	pattern *regexp.Regexp
}

// NewRegexAutomaton creates a regex automaton
func NewRegexAutomaton(pattern string) (*RegexAutomaton, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	
	return &RegexAutomaton{
		pattern: re,
	}, nil
}

// Accept tests if the automaton accepts the input
func (ra *RegexAutomaton) Accept(input []byte) bool {
	return ra.pattern.Match(input)
}

// RegexSearch performs regex search on the FSA
func RegexSearch(fsa FSA, pattern string) ([]string, error) {
	matcher, err := NewRegexMatcher(pattern)
	if err != nil {
		return nil, err
	}
	
	var results []string
	
	// Simple approach: test each key against the regex
	iter := fsa.Iterator()
	for iter.Next() {
		key := string(iter.Key())
		if matcher.pattern.MatchString(key) {
			results = append(results, key)
		}
	}
	
	return results, nil
}

// PrefixRegexSearch performs regex search on keys with a given prefix
func PrefixRegexSearch(fsa FSA, prefix, pattern string) ([]string, error) {
	matcher, err := NewRegexMatcher(pattern)
	if err != nil {
		return nil, err
	}
	
	var results []string
	
	// Use prefix iterator to find keys with the prefix
	iter := fsa.PrefixIterator([]byte(prefix))
	for iter.Next() {
		key := string(iter.Key())
		if matcher.pattern.MatchString(key) {
			results = append(results, key)
		}
	}
	
	return results, nil
}

// ComplexQuery represents a complex query combining multiple search types
type ComplexQuery struct {
	fsa FSA
}

// NewComplexQuery creates a new complex query
func NewComplexQuery(fsa FSA) *ComplexQuery {
	return &ComplexQuery{fsa: fsa}
}

// QueryResult represents the result of a complex query
type QueryResult struct {
	Keys  []string
	Count int
}

// Execute executes a complex query with multiple criteria
func (cq *ComplexQuery) Execute(options QueryOptions) (*QueryResult, error) {
	var candidates []string
	
	// Start with all keys or apply prefix filter
	if options.Prefix != "" {
		iterator := cq.fsa.PrefixIterator([]byte(options.Prefix))
		for iterator.Next() {
			key := iterator.Key()
			candidates = append(candidates, string(key))
		}
	} else if options.StartKey != "" || options.EndKey != "" {
		iterator := cq.fsa.RangeIterator([]byte(options.StartKey), []byte(options.EndKey))
		for iterator.Next() {
			key := iterator.Key()
			candidates = append(candidates, string(key))
		}
	} else {
		iterator := cq.fsa.Iterator()
		for iterator.Next() {
			key := iterator.Key()
			candidates = append(candidates, string(key))
		}
	}
	
	// Apply regex filter if specified
	if options.RegexPattern != "" {
		matcher, err := NewRegexMatcher(options.RegexPattern)
		if err != nil {
			return nil, err
		}
		
		filtered := make([]string, 0)
		for _, candidate := range candidates {
			if matcher.pattern.MatchString(candidate) {
				filtered = append(filtered, candidate)
			}
		}
		candidates = filtered
	}
	
	// Apply fuzzy search if specified
	if options.FuzzyPattern != "" {
		fuzzyResults := FuzzySearch(cq.fsa, options.FuzzyPattern, options.FuzzyMaxDistance)
		
		// Intersect with candidates
		candidateSet := make(map[string]bool)
		for _, candidate := range candidates {
			candidateSet[candidate] = true
		}
		
		filtered := make([]string, 0)
		for _, fuzzyResult := range fuzzyResults {
			if candidateSet[fuzzyResult] {
				filtered = append(filtered, fuzzyResult)
			}
		}
		candidates = filtered
	}
	
	// Apply limit if specified
	if options.Limit > 0 && len(candidates) > options.Limit {
		candidates = candidates[:options.Limit]
	}
	
	return &QueryResult{
		Keys:  candidates,
		Count: len(candidates),
	}, nil
}

// QueryOptions represents options for complex queries
type QueryOptions struct {
	Prefix            string // Prefix filter
	StartKey          string // Range start (inclusive)
	EndKey            string // Range end (exclusive)  
	RegexPattern      string // Regex pattern to match
	FuzzyPattern      string // Pattern for fuzzy search
	FuzzyMaxDistance  int    // Maximum edit distance for fuzzy search
	Limit             int    // Maximum number of results (0 = no limit)
}
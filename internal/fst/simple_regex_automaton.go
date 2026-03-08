package fst

import (
	"regexp"
	"sort"
	"strings"
)

// SimpleRegexAutomaton provides a simplified regex automaton that uses Go's regexp internally
// This is a fallback implementation while we debug the full NFA construction
type SimpleRegexAutomaton struct {
	pattern string
	regex   *regexp.Regexp
}

// NewSimpleRegexAutomaton creates a simple regex automaton
func NewSimpleRegexAutomaton(pattern string) (*SimpleRegexAutomaton, error) {
	// For FST intersection, anchor the pattern to match complete keys
	anchoredPattern := pattern
	if !strings.HasPrefix(pattern, "^") && !strings.HasSuffix(pattern, "$") {
		anchoredPattern = "^" + pattern + "$"
	}

	regex, err := regexp.Compile(anchoredPattern)
	if err != nil {
		return nil, err
	}

	return &SimpleRegexAutomaton{
		pattern: pattern,
		regex:   regex,
	}, nil
}

// TrueAutomataIntersection performs intersection using Go's regex engine
func (sra *SimpleRegexAutomaton) TrueAutomataIntersection(fst *FST) ([]string, error) {
	results := make([]string, 0)
	
	// Iterate through all FST keys and test each one
	iter := fst.Iterator()
	for iter.HasNext() {
		key, _ := iter.Next()
		keyStr := string(key)
		
		if sra.regex.MatchString(keyStr) {
			results = append(results, keyStr)
		}
	}
	
	sort.Strings(results)
	return results, nil
}

// MatchString tests if a string matches the regex
func (sra *SimpleRegexAutomaton) MatchString(s string) bool {
	return sra.regex.MatchString(s)
}

// IntersectWithFST provides the same interface as TrueRegexAutomaton
func (sra *SimpleRegexAutomaton) IntersectWithFST(fsa FSA) ([]string, error) {
	results := make([]string, 0)
	
	iter := fsa.Iterator()
	for iter.Next() {
		key := string(iter.Key())
		if sra.regex.MatchString(key) {
			results = append(results, key)
		}
	}
	
	sort.Strings(results)
	return results, nil
}
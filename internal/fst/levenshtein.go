package fst

import (
	"sort"
)

// LevenshteinState represents a state in the Levenshtein automaton
type LevenshteinState struct {
	Position int  // Position in the target string
	Errors   int  // Number of errors so far
	IsValid  bool // Whether this state is reachable
}

// LevenshteinAutomaton represents an automaton for fuzzy string matching
// using edit distance (insertions, deletions, substitutions)
type LevenshteinAutomaton struct {
	Pattern     string
	MaxDistance int
	States      [][]LevenshteinState // [position][errors] -> state
}

// NewLevenshteinAutomaton creates a Levenshtein automaton for fuzzy matching
func NewLevenshteinAutomaton(pattern string, maxDistance int) *LevenshteinAutomaton {
	patternLen := len(pattern)
	
	// Create state table: [position][errors]
	states := make([][]LevenshteinState, patternLen+maxDistance+1)
	for i := range states {
		states[i] = make([]LevenshteinState, maxDistance+1)
	}
	
	// Initialize starting states
	for e := 0; e <= maxDistance; e++ {
		states[e][e] = LevenshteinState{
			Position: e,
			Errors:   e,
			IsValid:  true,
		}
	}
	
	return &LevenshteinAutomaton{
		Pattern:     pattern,
		MaxDistance: maxDistance,
		States:      states,
	}
}

// Step advances the automaton with the given character
func (la *LevenshteinAutomaton) Step(char byte) *LevenshteinAutomaton {
	patternLen := len(la.Pattern)
	newStates := make([][]LevenshteinState, patternLen+la.MaxDistance+1)
	for i := range newStates {
		newStates[i] = make([]LevenshteinState, la.MaxDistance+1)
	}
	
	// For each current state, compute possible next states
	for pos := 0; pos < len(la.States); pos++ {
		for err := 0; err <= la.MaxDistance; err++ {
			currentState := la.States[pos][err]
			if !currentState.IsValid {
				continue
			}
			
			// Match transition (no error if characters match)
			if pos < patternLen {
				nextPos := pos + 1
				nextErr := err
				if la.Pattern[pos] != char {
					nextErr++
				}
				
				if nextErr <= la.MaxDistance && nextPos < len(newStates) {
					newStates[nextPos][nextErr] = LevenshteinState{
						Position: nextPos,
						Errors:   nextErr,
						IsValid:  true,
					}
				}
			}
			
			// Insertion (advance input, don't advance pattern)
			if err+1 <= la.MaxDistance && pos < len(newStates) {
				newStates[pos][err+1] = LevenshteinState{
					Position: pos,
					Errors:   err + 1,
					IsValid:  true,
				}
			}
			
			// Deletion (advance pattern, don't advance input)
			if pos < patternLen && err+1 <= la.MaxDistance && pos+1 < len(newStates) {
				newStates[pos+1][err+1] = LevenshteinState{
					Position: pos + 1,
					Errors:   err + 1,
					IsValid:  true,
				}
			}
		}
	}
	
	return &LevenshteinAutomaton{
		Pattern:     la.Pattern,
		MaxDistance: la.MaxDistance,
		States:      newStates,
	}
}

// IsMatch checks if the current state represents a successful match
func (la *LevenshteinAutomaton) IsMatch() bool {
	patternLen := len(la.Pattern)
	
	// Check if we can reach the end of the pattern within max distance
	for err := 0; err <= la.MaxDistance; err++ {
		// Direct match at end of pattern
		if patternLen < len(la.States) && la.States[patternLen][err].IsValid {
			return true
		}
		
		// Allow for trailing deletions (extra characters in pattern)
		for pos := patternLen; pos < len(la.States) && pos <= patternLen+err; pos++ {
			if la.States[pos][err].IsValid {
				return true
			}
		}
	}
	
	return false
}

// CanMatch checks if this automaton could potentially match with more input
func (la *LevenshteinAutomaton) CanMatch() bool {
	// Check if any state is still valid
	for i := range la.States {
		for j := range la.States[i] {
			if la.States[i][j].IsValid {
				return true
			}
		}
	}
	return false
}

// FuzzySearch performs fuzzy search on the FSA using Levenshtein distance
func FuzzySearch(fsa FSA, pattern string, maxDistance int) []string {
	var results []string
	
	// For SimpleFSA, we can access keys directly
	if simpleFSA, ok := fsa.(*SimpleFSA); ok {
		automaton := NewLevenshteinAutomaton(pattern, maxDistance)
		fuzzySearchRecursive(simpleFSA.keys, "", 0, automaton, &results)
	} else {
		// For other FSA implementations, iterate through all keys
		iter := fsa.Iterator()
		for iter.Next() {
			key := string(iter.Key())
			distance := computeLevenshteinDistance(key, pattern)
			if distance <= maxDistance {
				results = append(results, key)
			}
		}
	}
	
	sort.Strings(results)
	return results
}

// fuzzySearchRecursive performs recursive fuzzy search
func fuzzySearchRecursive(keys [][]byte, prefix string, index int, automaton *LevenshteinAutomaton, results *[]string) {
	// Check if current state is a match
	if index < len(keys) && string(keys[index]) == prefix && automaton.IsMatch() {
		*results = append(*results, prefix)
	}
	
	// If automaton can't match anymore, prune this branch
	if !automaton.CanMatch() {
		return
	}
	
	// Try all possible next characters
	tried := make(map[byte]bool)
	
	// Look at keys that have this prefix
	for i := index; i < len(keys); i++ {
		key := string(keys[i])
		
		// If key doesn't start with current prefix, we're done with this branch
		if len(key) <= len(prefix) || !hasPrefix(key, prefix) {
			if len(prefix) > 0 && !hasPrefix(key, prefix[:len(prefix)-1]) {
				break
			}
			continue
		}
		
		nextChar := key[len(prefix)]
		if tried[nextChar] {
			continue
		}
		tried[nextChar] = true
		
		// Step the automaton with this character
		nextAutomaton := automaton.Step(nextChar)
		if nextAutomaton.CanMatch() {
			fuzzySearchRecursive(keys, prefix+string(nextChar), i, nextAutomaton, results)
		}
	}
}

// Helper function to check if string has prefix
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// computeLevenshteinDistance computes the Levenshtein distance between two strings
func computeLevenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)
	
	// Create a matrix for dynamic programming
	dp := make([][]int, len1+1)
	for i := range dp {
		dp[i] = make([]int, len2+1)
	}
	
	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		dp[0][j] = j
	}
	
	// Fill the matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			
			dp[i][j] = min(
				dp[i-1][j]+1,      // deletion
				dp[i][j-1]+1,      // insertion
				dp[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return dp[len1][len2]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
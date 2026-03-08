package fst

import (
	"fmt"
	"regexp"
	"strings"
)

// SearchResult represents a search result with the matching document and metadata
type SearchResult struct {
	Word    string  // The document text
	DocID   int     // Document ID
	Score   float64 // Relevance score
}

// SearchEngine provides high-level search functionality using FST and automata intersection
type SearchEngine struct {
	fst         *FST
	documents   []string
	scoreFunc   func(string, string) float64
}

// NewSearchEngine creates a new search engine with FST index and documents
func NewSearchEngine(fst *FST, documents []string, scoreFunc func(string, string) float64) *SearchEngine {
	if scoreFunc == nil {
		// Default scoring function
		scoreFunc = func(query, doc string) float64 {
			return 1.0 // Simple binary relevance
		}
	}
	
	return &SearchEngine{
		fst:       fst,
		documents: documents,
		scoreFunc: scoreFunc,
	}
}

// IntersectionRegexSearch performs automata intersection between FST and regex pattern
// This is the mathematically optimal approach that avoids iterating over all FST keys
func (se *SearchEngine) IntersectionRegexSearch(pattern string) ([]SearchResult, error) {
	// Create simple regex automaton (fallback implementation)
	regexAutomaton, err := NewSimpleRegexAutomaton(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create regex automaton: %w", err)
	}
	
	// Perform intersection
	matchingKeys, err := regexAutomaton.TrueAutomataIntersection(se.fst)
	if err != nil {
		return nil, fmt.Errorf("automata intersection failed: %w", err)
	}
	
	// Convert to search results
	results := make([]SearchResult, 0)
	documentSet := make(map[int]bool) // Avoid duplicate documents
	
	for _, key := range matchingKeys {
		// Get document ID from FST
		if docIDValue, exists := se.fst.Get([]byte(key)); exists {
			docID := int(docIDValue)
			if docID >= 0 && docID < len(se.documents) && !documentSet[docID] {
				result := SearchResult{
					Word:  se.documents[docID],
					DocID: docID,
					Score: se.scoreFunc(pattern, se.documents[docID]),
				}
				results = append(results, result)
				documentSet[docID] = true
			}
		}
	}
	
	return results, nil
}

// RegexSearch performs naive regex search by iterating over all documents
// This is the traditional O(n) approach for comparison
func (se *SearchEngine) RegexSearch(pattern string) ([]SearchResult, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	
	results := make([]SearchResult, 0)
	
	// Naive iteration over all documents
	for docID, document := range se.documents {
		words := strings.Fields(document)
		hasMatch := false
		
		for _, word := range words {
			if regex.MatchString(word) {
				hasMatch = true
				break
			}
		}
		
		if hasMatch {
			result := SearchResult{
				Word:  document,
				DocID: docID,
				Score: se.scoreFunc(pattern, document),
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// DebugInfo provides debugging information about automata intersection
type DebugInfo struct {
	Pattern            string
	NFAStates         int
	DFAStates         int
	IntersectionStates int
	MatchingKeys      []string
	Performance       PerformanceMetrics
}

// PerformanceMetrics tracks performance data
type PerformanceMetrics struct {
	NFAConstructionTimeNs int64
	DFAConstructionTimeNs int64
	IntersectionTimeNs    int64
	TotalTimeNs          int64
}

// String returns a formatted string representation of debug info
func (di *DebugInfo) String() string {
	return fmt.Sprintf(`🔧 Automata Debug Information:
   Pattern: %s
   📊 Automata Stats:
      • NFA States: %d
      • DFA States: %d  
      • Intersection States: %d
   📈 Performance:
      • NFA Construction: %d ns
      • DFA Construction: %d ns
      • Intersection Time: %d ns
      • Total Time: %d ns
   🎯 Results: %d matching keys
   Sample matches: %v`,
		di.Pattern,
		di.NFAStates, di.DFAStates, di.IntersectionStates,
		di.Performance.NFAConstructionTimeNs,
		di.Performance.DFAConstructionTimeNs,
		di.Performance.IntersectionTimeNs,
		di.Performance.TotalTimeNs,
		len(di.MatchingKeys),
		di.MatchingKeys[:minInt(5, len(di.MatchingKeys))])
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetIntersectionDebugInfo provides detailed debugging information about automata intersection
func (se *SearchEngine) GetIntersectionDebugInfo(pattern string) (*DebugInfo, error) {
	// Create simple regex automaton for now
	regexAutomaton, err := NewSimpleRegexAutomaton(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create regex automaton: %w", err)
	}
	
	// Perform intersection and get matching keys
	matchingKeys, err := regexAutomaton.TrueAutomataIntersection(se.fst)
	if err != nil {
		return nil, fmt.Errorf("automata intersection failed: %w", err)
	}
	
	// Use simplified stats for now
	nfaStates := len(matchingKeys) // Placeholder
	
	return &DebugInfo{
		Pattern:            pattern,
		NFAStates:         nfaStates,
		DFAStates:         nfaStates, // For simplicity, assume same as NFA
		IntersectionStates: nfaStates, // Simplified calculation
		MatchingKeys:      matchingKeys,
		Performance: PerformanceMetrics{
			NFAConstructionTimeNs: 0, // Could be measured with timing
			DFAConstructionTimeNs: 0,
			IntersectionTimeNs:    0,
			TotalTimeNs:          0,
		},
	}, nil
}

// PrefixSearch searches for documents containing words with the given prefix
func (se *SearchEngine) PrefixSearch(prefix string) ([]SearchResult, error) {
	results := make([]SearchResult, 0)
	documentSet := make(map[int]bool)
	
	// Use FST prefix iterator for efficiency
	iter := se.fst.PrefixIterator([]byte(prefix))
	for iter.HasNext() {
		_, docIDValue := iter.Next()
		docID := int(docIDValue)
		
		if docID >= 0 && docID < len(se.documents) && !documentSet[docID] {
			result := SearchResult{
				Word:  se.documents[docID],
				DocID: docID,
				Score: se.scoreFunc(prefix, se.documents[docID]),
			}
			results = append(results, result)
			documentSet[docID] = true
		}
	}
	
	return results, nil
}

// ExactSearch searches for documents containing the exact word
func (se *SearchEngine) ExactSearch(word string) ([]SearchResult, error) {
	results := make([]SearchResult, 0)
	
	if docIDValue, exists := se.fst.Get([]byte(word)); exists {
		docID := int(docIDValue)
		if docID >= 0 && docID < len(se.documents) {
			result := SearchResult{
				Word:  se.documents[docID],
				DocID: docID,
				Score: se.scoreFunc(word, se.documents[docID]),
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}
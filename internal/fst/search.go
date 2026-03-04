/*
 * Copyright (c) 2011 jamra.source@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package fst

import (
	"sort"
	"strings"

	"github.com/jamra/gocleo/internal/search"
	"github.com/jamra/gocleo/internal/scoring"
)

// SearchEngine provides search functionality using FST
type SearchEngine struct {
	fst         *FST
	documents   []string // Document content indexed by document ID
	scoringFunc scoring.ScoringFunction
}

// NewSearchEngine creates a new FST-based search engine
func NewSearchEngine(fst *FST, documents []string, scoringFunc scoring.ScoringFunction) *SearchEngine {
	if scoringFunc == nil {
		scoringFunc = scoring.DefaultScore
	}

	return &SearchEngine{
		fst:         fst,
		documents:   documents,
		scoringFunc: scoringFunc,
	}
}

// Search performs a search using the FST and returns ranked results
func (e *SearchEngine) Search(query string) []search.RankedResult {
	if query == "" || e.fst == nil {
		return []search.RankedResult{}
	}

	results := make([]search.RankedResult, 0)
	queryLower := strings.ToLower(query)

	// Use prefix iterator to find all matches
	iter := e.fst.PrefixIterator([]byte(queryLower))
	
	for iter.HasNext() {
		_, docID := iter.Next()
		
		// Ensure we have a valid document ID
		if int(docID) <= len(e.documents) && docID > 0 {
			docContent := e.documents[docID-1] // Assuming 1-based document IDs
			
			// Score the match
			score := e.scoringFunc(query, docContent)
			
			if score > 0 {
				results = append(results, search.RankedResult{
					Word:  docContent,
					Score: score,
				})
			}
		}
	}

	// Sort results by score (descending)
	sort.Sort(search.ByScore{search.RankedResults(results)})
	
	return results
}

// ExactSearch performs exact matching using the FST
func (e *SearchEngine) ExactSearch(query string) []search.RankedResult {
	if query == "" || e.fst == nil {
		return []search.RankedResult{}
	}

	queryLower := strings.ToLower(query)
	
	if docID, found := e.fst.Get([]byte(queryLower)); found {
		if int(docID) <= len(e.documents) && docID > 0 {
			docContent := e.documents[docID-1]
			score := e.scoringFunc(query, docContent)
			
			if score > 0 {
				return []search.RankedResult{
					{
						Word:  docContent,
						Score: score,
					},
				}
			}
		}
	}
	
	return []search.RankedResult{}
}

// FuzzySearch performs fuzzy matching by expanding the search to similar terms
func (e *SearchEngine) FuzzySearch(query string, maxDistance int) []search.RankedResult {
	if query == "" || e.fst == nil {
		return []search.RankedResult{}
	}

	results := make([]search.RankedResult, 0)
	queryLower := strings.ToLower(query)

	// Iterate through all keys and calculate edit distance
	iter := e.fst.Iterator()
	
	for iter.HasNext() {
		key, docID := iter.Next()
		keyStr := string(key)
		
		// Calculate Levenshtein distance using the existing function in this package
		distance := computeLevenshteinDistance(queryLower, keyStr)
		
		if distance <= maxDistance {
			if int(docID) <= len(e.documents) && docID > 0 {
				docContent := e.documents[docID-1]
				
				// Adjust score based on edit distance
				baseScore := e.scoringFunc(query, docContent)
				adjustedScore := baseScore * (1.0 - float64(distance)/float64(len(queryLower)+len(keyStr)))
				
				if adjustedScore > 0 {
					results = append(results, search.RankedResult{
						Word:  docContent,
						Score: adjustedScore,
					})
				}
			}
		}
	}

	// Sort results by score (descending)
	sort.Sort(search.ByScore{search.RankedResults(results)})
	
	return results
}

// SetScoringFunction updates the scoring function
func (e *SearchEngine) SetScoringFunction(scoringFunc scoring.ScoringFunction) {
	if scoringFunc != nil {
		e.scoringFunc = scoringFunc
	}
}

// GetStats returns statistics about the FST search engine
func (e *SearchEngine) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"fst_size":       e.fst.Size(),
		"fst_empty":      e.fst.IsEmpty(),
		"documents":      len(e.documents),
	}
}

// BuildFSTFromDocuments creates an FST from a slice of documents
func BuildFSTFromDocuments(documents []string) (*FST, error) {
	if len(documents) == 0 {
		return NewFSTBuilder().Build()
	}

	// Create a map to collect unique words and their document IDs
	wordToID := make(map[string]uint64)
	
	for docID, doc := range documents {
		// Simple tokenization - split by whitespace and punctuation
		words := tokenize(strings.ToLower(doc))
		
		for _, word := range words {
			if word != "" {
				// Use the first document ID where this word appears
				if _, exists := wordToID[word]; !exists {
					wordToID[word] = uint64(docID + 1) // 1-based document IDs
				}
			}
		}
	}

	// Convert to sorted slice for FST building
	words := make([]string, 0, len(wordToID))
	for word := range wordToID {
		words = append(words, word)
	}
	sort.Strings(words)

	// Build the FST
	builder := NewFSTBuilder()
	for _, word := range words {
		docID := wordToID[word]
		err := builder.Add([]byte(word), docID)
		if err != nil {
			return nil, err
		}
	}

	return builder.Build()
}

// tokenize splits text into words
func tokenize(text string) []string {
	// Simple tokenization - this could be enhanced
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
	})
	
	return words
}

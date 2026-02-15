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

package search

import (
	"github.com/jamra/gocleo/internal/bloom"
	"github.com/jamra/gocleo/internal/index"
	"github.com/jamra/gocleo/internal/scoring"
)

// Engine provides the core search functionality.
type Engine struct {
	invertedIndex *index.InvertedIndex
	forwardIndex  *index.ForwardIndex
	scoringFunc   scoring.ScoringFunction
}

// NewEngine creates a new search engine with the provided indexes and scoring function.
func NewEngine(invertedIndex *index.InvertedIndex, forwardIndex *index.ForwardIndex, scoringFunc scoring.ScoringFunction) *Engine {
	if scoringFunc == nil {
		scoringFunc = scoring.DefaultScore
	}

	return &Engine{
		invertedIndex: invertedIndex,
		forwardIndex:  forwardIndex,
		scoringFunc:   scoringFunc,
	}
}

// Search performs a Cleo search query and returns ranked results.
// The search process:
// 1. Get candidates from the inverted index based on query prefix
// 2. Filter candidates using bloom filter matching
// 3. Score and rank the filtered candidates
func (e *Engine) Search(query string) []RankedResult {
	if query == "" {
		return []RankedResult{}
	}

	// Step 1: Get candidates from inverted index
	candidates := e.invertedIndex.Search(query)
	if candidates == nil {
		return []RankedResult{}
	}

	// Step 2: Filter using bloom filters and score
	results := make([]RankedResult, 0)
	queryBloom := bloom.ComputeBloomFilter(query)

	for _, candidate := range candidates {
		// Test bloom filter match
		if bloom.TestBytesFromQuery(candidate.Score, queryBloom) {
			// Get the actual document content from forward index
			docContent := e.forwardIndex.ItemAt(candidate.Id)
			
			// Score the match
			score := e.scoringFunc(query, docContent)
			
			if score > 0 { // Only include results with positive scores
				results = append(results, RankedResult{
					Word:  docContent,
					Score: score,
				})
			}
		}
	}

	return results
}

// SetScoringFunction updates the scoring function used by the search engine.
func (e *Engine) SetScoringFunction(scoringFunc scoring.ScoringFunction) {
	if scoringFunc != nil {
		e.scoringFunc = scoringFunc
	}
}

// GetIndexStats returns statistics about the search indexes.
func (e *Engine) GetIndexStats() map[string]interface{} {
	return map[string]interface{}{
		"inverted_index_prefixes": e.invertedIndex.Size(),
		"inverted_index_documents": e.invertedIndex.GetDocumentCount(),
		"forward_index_documents": e.forwardIndex.Size(),
	}
}

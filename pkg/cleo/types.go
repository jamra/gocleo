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

// Package cleo provides a fast prefix search algorithm optimized for large text corpora.
// This is the public API for the Cleo search library.
package cleo

import (
	"github.com/jamra/gocleo/internal/scoring"
	"github.com/jamra/gocleo/internal/search"
)

// Result represents a search result with its relevance score.
type Result struct {
	Word  string  `json:"word"`  // The matched word
	Score float64 `json:"score"` // The relevance score (0-1, higher is better)
}

// ScoringFunction defines the interface for custom scoring functions.
// It takes a query and candidate word, returning a relevance score.
type ScoringFunction = scoring.ScoringFunction

// Config holds configuration options for a Cleo search instance.
type Config struct {
	// ScoringFunction defines how to score matches. If nil, uses DefaultScore.
	ScoringFunction ScoringFunction
	
	// MaxResults limits the number of results returned. 0 means no limit.
	MaxResults int
	
	// MinScore filters out results below this threshold. 0 means no filtering.
	MinScore float64
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScoringFunction: nil, // Will use scoring.DefaultScore
		MaxResults:      0,   // No limit
		MinScore:        0.0, // No filtering
	}
}

// convertResults converts internal search results to public API results.
func convertResults(results []search.RankedResult) []Result {
	apiResults := make([]Result, len(results))
	for i, result := range results {
		apiResults[i] = Result{
			Word:  result.Word,
			Score: result.Score,
		}
	}
	return apiResults
}

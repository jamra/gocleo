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

package scoring

import (
	"math"
	"strings"
)

// ScoringFunction defines the signature for scoring functions.
type ScoringFunction func(word, query string) float64

// DefaultScore computes the default Cleo score using Levenshtein distance
// and Jaccard coefficient normalization.
func DefaultScore(query, candidate string) float64 {
	// Normalize inputs to lowercase for case-insensitive matching
	queryLower := strings.ToLower(query)
	candidateLower := strings.ToLower(candidate)

	// Calculate Levenshtein distance
	levDist := float64(LevenshteinDistance(queryLower, candidateLower))

	// Calculate Jaccard coefficient for normalization
	jaccard := JaccardCoefficient(queryLower, candidateLower)

	// Avoid division by zero
	if jaccard == 0 {
		return 0
	}

	// Normalize the distance by the Jaccard coefficient
	// Lower distances with higher Jaccard coefficients get higher scores
	score := (1.0 / (1.0 + levDist)) * jaccard
	return score
}

// JaccardCoefficient computes the Jaccard coefficient between two strings.
// This measures the similarity between the character sets of the strings.
func JaccardCoefficient(s1, s2 string) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Create character sets
	set1 := make(map[rune]bool)
	set2 := make(map[rune]bool)

	for _, char := range s1 {
		set1[char] = true
	}
	for _, char := range s2 {
		set2[char] = true
	}

	// Calculate intersection and union
	intersection := 0
	for char := range set1 {
		if set2[char] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// PrefixScore gives higher scores to candidates that start with the query.
func PrefixScore(query, candidate string) float64 {
	queryLower := strings.ToLower(query)
	candidateLower := strings.ToLower(candidate)

	if strings.HasPrefix(candidateLower, queryLower) {
		// Perfect prefix match gets higher score
		return 1.0 - (float64(len(candidate)-len(query)) / float64(len(candidate)))
	}

	// Fall back to default scoring
	return DefaultScore(query, candidate) * 0.5
}

// ExactScore prioritizes exact matches and close prefixes.
func ExactScore(query, candidate string) float64 {
	queryLower := strings.ToLower(query)
	candidateLower := strings.ToLower(candidate)

	if queryLower == candidateLower {
		return 1.0 // Exact match
	}

	if strings.HasPrefix(candidateLower, queryLower) {
		return 0.9 // Prefix match
	}

	if strings.Contains(candidateLower, queryLower) {
		return 0.7 // Contains match
	}

	// Use default scoring for other cases
	return DefaultScore(query, candidate)
}

// FuzzyScore emphasizes fuzzy matching using only Levenshtein distance.
func FuzzyScore(query, candidate string) float64 {
	queryLower := strings.ToLower(query)
	candidateLower := strings.ToLower(candidate)

	levDist := float64(LevenshteinDistance(queryLower, candidateLower))
	maxLen := math.Max(float64(len(query)), float64(len(candidate)))

	if maxLen == 0 {
		return 1.0
	}

	// Normalize by maximum length
	return 1.0 - (levDist / maxLen)
}

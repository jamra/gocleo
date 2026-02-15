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

// Package search provides the core search functionality for Cleo.
package search

// RankedResult represents a search result with its score.
type RankedResult struct {
	Word  string  `json:"word"`  // The matched word/document
	Score float64 `json:"score"` // The relevance score
}

// RankedResults is a slice of RankedResult for sorting.
type RankedResults []RankedResult

// Len implements sort.Interface
func (r RankedResults) Len() int { return len(r) }

// Swap implements sort.Interface  
func (r RankedResults) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// ByScore implements sort.Interface for sorting by score (descending).
type ByScore struct{ RankedResults }

// Less implements sort.Interface for descending score order
func (s ByScore) Less(i, j int) bool { 
	return s.RankedResults[i].Score > s.RankedResults[j].Score 
}

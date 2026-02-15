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

package cleo

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/jamra/gocleo/internal/bloom"
	"github.com/jamra/gocleo/internal/index"
	"github.com/jamra/gocleo/internal/scoring"
	"github.com/jamra/gocleo/internal/search"
)

// Client represents a Cleo search client instance.
// It's thread-safe and can be used concurrently.
type Client struct {
	engine *search.Engine
	config *Config
	mutex  sync.RWMutex
}

// New creates a new Cleo search client from a corpus file.
// The corpusPath should point to a newline-separated text file where each line is a searchable term.
func New(corpusPath string, config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create indexes
	invertedIndex := index.NewInvertedIndex()
	forwardIndex := index.NewForwardIndex()

	// Load the corpus
	err := loadCorpus(corpusPath, invertedIndex, forwardIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to load corpus: %w", err)
	}

	// Create the search engine
	engine := search.NewEngine(invertedIndex, forwardIndex, config.ScoringFunction)

	return &Client{
		engine: engine,
		config: config,
	}, nil
}

// NewFromWords creates a new Cleo search client from a slice of words.
// This is useful for programmatically creating search indexes.
func NewFromWords(words []string, config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create indexes
	invertedIndex := index.NewInvertedIndex()
	forwardIndex := index.NewForwardIndex()

	// Load words into indexes
	for docID, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		bloomFilter := bloom.ComputeBloomFilter(word)
		invertedIndex.AddDoc(docID+1, word, bloomFilter)
		forwardIndex.AddDoc(docID+1, word)
	}

	// Create the search engine
	engine := search.NewEngine(invertedIndex, forwardIndex, config.ScoringFunction)

	return &Client{
		engine: engine,
		config: config,
	}, nil
}

// Search performs a search query and returns ranked results.
func (c *Client) Search(query string) ([]Result, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if query == "" {
		return []Result{}, nil
	}

	// Perform the search
	results := c.engine.Search(query)

	// Sort results by score (descending)
	sort.Sort(search.ByScore{RankedResults: results})

	// Apply filtering and limits
	filtered := c.filterResults(results)

	return convertResults(filtered), nil
}

// SetScoringFunction updates the scoring function used for search.
// This is thread-safe and will affect all subsequent searches.
func (c *Client) SetScoringFunction(scoringFunc ScoringFunction) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.config.ScoringFunction = scoringFunc
	c.engine.SetScoringFunction(scoringFunc)
}

// GetStats returns statistics about the search indexes.
func (c *Client) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.engine.GetIndexStats()
}

// filterResults applies MinScore and MaxResults filtering.
func (c *Client) filterResults(results []search.RankedResult) []search.RankedResult {
	filtered := make([]search.RankedResult, 0)

	for _, result := range results {
		// Apply minimum score filter
		if c.config.MinScore > 0 && result.Score < c.config.MinScore {
			continue
		}

		filtered = append(filtered, result)

		// Apply maximum results limit
		if c.config.MaxResults > 0 && len(filtered) >= c.config.MaxResults {
			break
		}
	}

	return filtered
}

// loadCorpus loads words from a corpus file into the indexes.
func loadCorpus(corpusPath string, invertedIndex *index.InvertedIndex, forwardIndex *index.ForwardIndex) error {
	file, err := os.Open(corpusPath)
	if err != nil {
		return fmt.Errorf("failed to open corpus file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	docID := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		// Compute bloom filter for the word
		bloomFilter := bloom.ComputeBloomFilter(line)

		// Add to both indexes
		invertedIndex.AddDoc(docID, line, bloomFilter)
		forwardIndex.AddDoc(docID, line)

		docID++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading corpus file: %w", err)
	}

	return nil
}

// Predefined scoring functions for convenience

// DefaultScore is the default Cleo scoring function using Levenshtein distance and Jaccard coefficient.
var DefaultScore ScoringFunction = scoring.DefaultScore

// PrefixScore gives higher scores to candidates that start with the query.
var PrefixScore ScoringFunction = scoring.PrefixScore  

// ExactScore prioritizes exact matches and close prefixes.
var ExactScore ScoringFunction = scoring.ExactScore

// FuzzyScore emphasizes fuzzy matching using Levenshtein distance.
var FuzzyScore ScoringFunction = scoring.FuzzyScore

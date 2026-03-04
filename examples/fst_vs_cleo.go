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

// Package main demonstrates FST vs Cleo search comparison
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jamra/gocleo/internal/fst"
	"github.com/jamra/gocleo/internal/search"
	"github.com/jamra/gocleo/internal/scoring"
	"github.com/jamra/gocleo/internal/index"
	"github.com/jamra/gocleo/internal/bloom"
)

func main() {
	fmt.Println("=== FST vs Cleo Search Comparison ===")

	// Sample dataset - technology terms
	documents := []string{
		"algorithm design and analysis",
		"artificial intelligence machine learning",
		"backend development with golang",
		"blockchain technology and cryptocurrency",
		"cloud computing infrastructure",
		"computer science fundamentals",
		"database management systems",
		"distributed systems architecture",
		"frontend web development",
		"golang programming language",
		"machine learning algorithms",
		"neural networks deep learning",
		"programming best practices",
		"software engineering principles",
		"web application development",
		"data structures and algorithms",
		"microservices architecture patterns",
		"devops and continuous integration",
		"cybersecurity best practices",
		"mobile application development",
	}

	fmt.Printf("\nDataset: %d documents\n", len(documents))

	// Build FST search engine
	fmt.Println("\nBuilding FST search engine...")
	start := time.Now()
	
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		log.Fatalf("Error building FST: %v", err)
	}
	
	fstEngine := fst.NewSearchEngine(fstIndex, documents, scoring.DefaultScore)
	fstBuildTime := time.Since(start)
	
	fmt.Printf("FST built in: %v\n", fstBuildTime)
	fmt.Printf("FST stats: %+v\n", fstEngine.GetStats())

	// Build Cleo search engine
	fmt.Println("\nBuilding Cleo search engine...")
	start = time.Now()
	
	// Create inverted index
	invertedIndex := index.NewInvertedIndex()
	forwardIndex := index.NewForwardIndex()
	
	// Add documents to indexes
	for i, doc := range documents {
		docID := i
		forwardIndex.AddDoc(docID, doc)
		
		// Tokenize and add to inverted index
		words := tokenize(doc)
		for _, word := range words {
			bloomFilter := bloom.ComputeBloomFilter(word)
			invertedIndex.AddDoc(docID, word, bloomFilter)
		}
	}
	
	cleoEngine := search.NewEngine(invertedIndex, forwardIndex, scoring.DefaultScore)
	cleoBuildTime := time.Since(start)
	
	fmt.Printf("Cleo built in: %v\n", cleoBuildTime)
	fmt.Printf("Cleo stats: %+v\n", cleoEngine.GetIndexStats())

	// Test queries
	queries := []string{
		"algorithm",
		"golang",
		"machine",
		"web",
		"data",
		"development",
		"artificial",
		"blockchain",
	}

	fmt.Println("\n=== Search Performance Comparison ===")

	for _, query := range queries {
		fmt.Printf("\n--- Query: '%s' ---\n", query)

		// FST Search
		start = time.Now()
		fstResults := fstEngine.Search(query)
		fstSearchTime := time.Since(start)

		// Cleo Search
		start = time.Now()
		cleoResults := cleoEngine.Search(query)
		cleoSearchTime := time.Since(start)

		// Display results
		fmt.Printf("FST Search Time: %v, Results: %d\n", fstSearchTime, len(fstResults))
		fmt.Printf("Cleo Search Time: %v, Results: %d\n", cleoSearchTime, len(cleoResults))

		// Show top results
		fmt.Println("FST Top Results:")
		for i, result := range fstResults {
			if i >= 3 { break } // Show top 3
			fmt.Printf("  %.3f: %s\n", result.Score, result.Word)
		}

		fmt.Println("Cleo Top Results:")
		for i, result := range cleoResults {
			if i >= 3 { break } // Show top 3
			fmt.Printf("  %.3f: %s\n", result.Score, result.Word)
		}

		// Performance comparison
		if fstSearchTime < cleoSearchTime {
			speedup := float64(cleoSearchTime) / float64(fstSearchTime)
			fmt.Printf("FST is %.2fx faster\n", speedup)
		} else {
			speedup := float64(fstSearchTime) / float64(cleoSearchTime)
			fmt.Printf("Cleo is %.2fx faster\n", speedup)
		}
	}

	// Test fuzzy search (FST specific)
	fmt.Println("\n=== FST Fuzzy Search ===")
	fuzzyQueries := []string{"algorthm", "machin", "developmnt"}
	
	for _, query := range fuzzyQueries {
		fmt.Printf("\nFuzzy search for '%s' (max distance: 2):\n", query)
		start = time.Now()
		fuzzyResults := fstEngine.FuzzySearch(query, 2)
		fuzzySearchTime := time.Since(start)
		
		fmt.Printf("Time: %v, Results: %d\n", fuzzySearchTime, len(fuzzyResults))
		for i, result := range fuzzyResults {
			if i >= 5 { break } // Show top 5
			fmt.Printf("  %.3f: %s\n", result.Score, result.Word)
		}
	}

	// Memory usage comparison
	fmt.Println("\n=== Summary ===")
	fmt.Printf("Build Time - FST: %v, Cleo: %v\n", fstBuildTime, cleoBuildTime)
	fmt.Printf("FST Features: Exact match, Prefix search, Fuzzy search, Range queries\n")
	fmt.Printf("Cleo Features: Bloom filter optimization, Inverted index, Prefix search\n")
	
	fmt.Println("\n=== FST vs Cleo Comparison Complete ===")
}

// tokenize splits text into words (same as in search.go)
func tokenize(text string) []string {
	words := make([]string, 0)
	currentWord := make([]rune, 0)
	
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			currentWord = append(currentWord, r)
		} else {
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = currentWord[:0]
			}
		}
	}
	
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}
	
	return words
}

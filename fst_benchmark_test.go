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
	"strings"
	"testing"

	"github.com/jamra/gocleo/internal/fst"
	"github.com/jamra/gocleo/internal/search"
	"github.com/jamra/gocleo/internal/scoring"
	"github.com/jamra/gocleo/internal/index"
	"github.com/jamra/gocleo/internal/bloom"
)

// generateTestData creates a larger dataset for benchmarking
func generateTestData(size int) []string {
	baseWords := []string{
		"algorithm", "analysis", "application", "architecture", "artificial",
		"backend", "blockchain", "browser", "build", "business",
		"cloud", "code", "computer", "container", "continuous",
		"data", "database", "deployment", "design", "development",
		"distributed", "documentation", "domain", "engineering", "framework",
		"frontend", "function", "golang", "infrastructure", "integration",
		"javascript", "kubernetes", "language", "learning", "machine",
		"microservice", "mobile", "network", "neural", "optimization",
		"performance", "programming", "protocol", "python", "quality",
		"security", "server", "software", "system", "technology",
		"testing", "tools", "user", "version", "web",
	}

	documents := make([]string, size)
	for i := 0; i < size; i++ {
		// Create realistic documents by combining 3-5 words
		doc := ""
		for j := 0; j < 3+i%3; j++ {
			if j > 0 {
				doc += " "
			}
			doc += baseWords[i%len(baseWords)] + baseWords[(i+j)%len(baseWords)]
		}
		documents[i] = doc
	}
	return documents
}

func setupFSTEngine(documents []string) *fst.SearchEngine {
	fstIndex, _ := fst.BuildFSTFromDocuments(documents)
	return fst.NewSearchEngine(fstIndex, documents, scoring.DefaultScore)
}

func setupCleoEngine(documents []string) *search.Engine {
	invertedIndex := index.NewInvertedIndex()
	forwardIndex := index.NewForwardIndex()
	
	for i, doc := range documents {
		docID := i
		forwardIndex.AddDoc(docID, doc)
		
		words := tokenizeForBench(doc)
		for _, word := range words {
			bloomFilter := bloom.ComputeBloomFilter(word)
			invertedIndex.AddDoc(docID, word, bloomFilter)
		}
	}
	
	return search.NewEngine(invertedIndex, forwardIndex, scoring.DefaultScore)
}

func tokenizeForBench(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// Benchmark FST build time
func BenchmarkFSTBuild100(b *testing.B) {
	documents := generateTestData(100)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fst.BuildFSTFromDocuments(documents)
	}
}

func BenchmarkFSTBuild1000(b *testing.B) {
	documents := generateTestData(1000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fst.BuildFSTFromDocuments(documents)
	}
}

func BenchmarkFSTBuild10000(b *testing.B) {
	documents := generateTestData(10000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fst.BuildFSTFromDocuments(documents)
	}
}

// Benchmark Cleo build time
func BenchmarkCleoBuild100(b *testing.B) {
	documents := generateTestData(100)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		setupCleoEngine(documents)
	}
}

func BenchmarkCleoBuild1000(b *testing.B) {
	documents := generateTestData(1000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		setupCleoEngine(documents)
	}
}

func BenchmarkCleoBuild10000(b *testing.B) {
	documents := generateTestData(10000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		setupCleoEngine(documents)
	}
}

// Benchmark FST search
func BenchmarkFSTSearch(b *testing.B) {
	documents := generateTestData(1000)
	engine := setupFSTEngine(documents)
	queries := []string{"algorithm", "data", "web", "machine", "development"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.Search(query)
	}
}

func BenchmarkFSTSearchLarge(b *testing.B) {
	documents := generateTestData(10000)
	engine := setupFSTEngine(documents)
	queries := []string{"algorithm", "data", "web", "machine", "development"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.Search(query)
	}
}

// Benchmark Cleo search
func BenchmarkCleoSearch(b *testing.B) {
	documents := generateTestData(1000)
	engine := setupCleoEngine(documents)
	queries := []string{"algorithm", "data", "web", "machine", "development"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.Search(query)
	}
}

func BenchmarkCleoSearchLarge(b *testing.B) {
	documents := generateTestData(10000)
	engine := setupCleoEngine(documents)
	queries := []string{"algorithm", "data", "web", "machine", "development"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.Search(query)
	}
}

// Benchmark FST fuzzy search
func BenchmarkFSTFuzzySearch(b *testing.B) {
	documents := generateTestData(1000)
	engine := setupFSTEngine(documents)
	queries := []string{"algorthm", "dat", "wb", "machin", "developmnt"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.FuzzySearch(query, 2)
	}
}

// Benchmark FST exact search
func BenchmarkFSTExactSearch(b *testing.B) {
	documents := generateTestData(1000)
	engine := setupFSTEngine(documents)
	queries := []string{"algorithm", "data", "web", "machine", "development"}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		engine.ExactSearch(query)
	}
}

// Comparative benchmark - run both engines on same query
func BenchmarkComparisonSearch(b *testing.B) {
	documents := generateTestData(1000)
	fstEngine := setupFSTEngine(documents)
	cleoEngine := setupCleoEngine(documents)
	query := "algorithm"
	
	b.Run("FST", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fstEngine.Search(query)
		}
	})
	
	b.Run("Cleo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cleoEngine.Search(query)
		}
	})
}

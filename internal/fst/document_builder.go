package fst

import (
	"sort"
	"strings"
)

// BuildFSTFromDocuments creates an FST from a collection of documents
// Each word in each document becomes a key with its document index as the value
func BuildFSTFromDocuments(documents []string) (*FST, error) {
	builder := NewFSTBuilder()
	wordToDocMap := make(map[string][]int)
	
	// Extract all unique words from all documents
	for docID, doc := range documents {
		words := extractWords(doc)
		for _, word := range words {
			if len(word) > 0 {
				wordToDocMap[word] = append(wordToDocMap[word], docID)
			}
		}
	}
	
	// Convert to sorted keys for FST building
	var keys []string
	for word := range wordToDocMap {
		keys = append(keys, word)
	}
	sort.Strings(keys)
	
	// Add to FST builder
	for _, word := range keys {
		// Use the first document ID as the value (for simplicity)
		// In a real implementation, you might want to store more complex data
		docID := wordToDocMap[word][0]
		err := builder.Add([]byte(word), uint64(docID))
		if err != nil {
			return nil, err
		}
	}
	
	return builder.Build()
}

// extractWords extracts individual words from a document
func extractWords(document string) []string {
	// Simple word extraction - split on whitespace and normalize
	words := strings.Fields(document)
	var result []string
	
	for _, word := range words {
		// Clean the word (remove punctuation, convert to lowercase)
		cleaned := strings.ToLower(strings.Trim(word, ".,!?;:()[]{}\"'"))
		if len(cleaned) > 0 {
			result = append(result, cleaned)
		}
	}
	
	return result
}

// BuildFSTFromWords creates an FST from a simple list of words
func BuildFSTFromWords(words []string) (*FST, error) {
	builder := NewFSTBuilder()
	
	// Sort words for FST building
	sortedWords := make([]string, len(words))
	copy(sortedWords, words)
	sort.Strings(sortedWords)
	
	// Remove duplicates and add to FST
	seen := make(map[string]bool)
	for i, word := range sortedWords {
		if !seen[word] && len(word) > 0 {
			err := builder.Add([]byte(word), uint64(i))
			if err != nil {
				return nil, err
			}
			seen[word] = true
		}
	}
	
	return builder.Build()
}
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

package index

// InvertedIndex maps query prefixes to matching documents with their bloom filters.
// This enables fast candidate retrieval for search queries.
type InvertedIndex map[string][]Document

// NewInvertedIndex creates a new empty inverted index.
func NewInvertedIndex() *InvertedIndex {
	i := make(InvertedIndex)
	return &i
}

// Size returns the number of prefixes in the index.
func (x *InvertedIndex) Size() int {
	return len(*x)
}

// AddDoc adds a document to the inverted index with the given document ID,
// content, and bloom filter score.
func (x *InvertedIndex) AddDoc(docId int, doc string, bloom int) {
	prefix := GetPrefix(doc)
	
	document := Document{
		Id:    docId,
		Score: bloom,
		Doc:   doc,
	}

	// Add to the index under the prefix key
	(*x)[prefix] = append((*x)[prefix], document)
}

// Search retrieves all documents that match the query prefix.
// Returns nil if no documents are found for the prefix.
func (x *InvertedIndex) Search(query string) []Document {
	prefix := GetPrefix(query)
	
	documents, found := (*x)[prefix]
	if found {
		return documents
	}
	return nil
}

// GetAllPrefixes returns all prefixes stored in the index.
func (x *InvertedIndex) GetAllPrefixes() []string {
	prefixes := make([]string, 0, len(*x))
	for prefix := range *x {
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}

// GetDocumentCount returns the total number of documents across all prefixes.
func (x *InvertedIndex) GetDocumentCount() int {
	count := 0
	for _, documents := range *x {
		count += len(documents)
	}
	return count
}

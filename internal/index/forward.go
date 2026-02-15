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

import "strings"

// ForwardIndex maps document IDs to their content.
// This enables fast document retrieval by ID during search result ranking.
type ForwardIndex map[int]string

// NewForwardIndex creates a new empty forward index.
func NewForwardIndex() *ForwardIndex {
	i := make(ForwardIndex)
	return &i
}

// AddDoc adds a document to the forward index.
// If the document ID already exists, it will be overwritten.
func (x *ForwardIndex) AddDoc(docId int, doc string) {
	// Store the first word from the document
	// This matches the original behavior which seems to index by word
	words := strings.Fields(doc)
	if len(words) > 0 {
		(*x)[docId] = words[0]
	} else {
		(*x)[docId] = doc
	}
}

// ItemAt retrieves the document content for the given document ID.
// Returns an empty string if the document ID is not found.
func (x *ForwardIndex) ItemAt(docId int) string {
	content, exists := (*x)[docId]
	if exists {
		return content
	}
	return ""
}

// Size returns the number of documents in the forward index.
func (x *ForwardIndex) Size() int {
	return len(*x)
}

// GetAllDocumentIds returns all document IDs in the index.
func (x *ForwardIndex) GetAllDocumentIds() []int {
	ids := make([]int, 0, len(*x))
	for id := range *x {
		ids = append(ids, id)
	}
	return ids
}

// Contains checks if a document ID exists in the index.
func (x *ForwardIndex) Contains(docId int) bool {
	_, exists := (*x)[docId]
	return exists
}

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

// Package index provides inverted and forward index implementations for Cleo search.
package index

// Document represents a document in the search index.
type Document struct {
	Id    int    `json:"id"`    // Document ID
	Score int    `json:"score"` // Bloom filter score
	Doc   string `json:"doc"`   // Document content
}

// GetPrefix extracts the search prefix from a query.
// Currently returns the first 4 characters or the entire string if shorter.
func GetPrefix(query string) string {
	if len(query) >= 4 {
		return query[0:4]
	}
	return query
}

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

// Package bloom provides bloom filter utilities for the Cleo search algorithm.
package bloom

import "fmt"

// ComputeBloomFilter computes the bloom filter for a given string.
// It uses a simple hash function to create a bloom filter representation
// that can be used for fast prefix matching.
func ComputeBloomFilter(s string) int {
	bloom := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		//first hash function
		h1 := (int(c) * 239) % 31

		//second hash function (reduces collisions for bloom)
		h2 := (int(c) * 991) % 31

		//create bit mask
		bloom = bloom | (1 << uint(h1))
		bloom = bloom | (1 << uint(h2))
	}
	return bloom
}

// TestBytesFromQuery tests if the bloom filter matches the query.
// It compares bits between the bloom filter (bf) and query bloom filter (qBloom).
func TestBytesFromQuery(bf int, qBloom int) bool {
	return (bf & qBloom) == qBloom
}

// DebugBloomFilter returns a string representation of the bloom filter for debugging.
func DebugBloomFilter(bloom int) string {
	return fmt.Sprintf("Bloom: %032b (%d)", bloom, bloom)
}

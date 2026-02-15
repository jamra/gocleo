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

// Package scoring provides string distance and scoring algorithms.
package scoring

// LevenshteinDistance computes the Levenshtein distance between two strings.
// This is the fixed implementation that handles edge cases properly.
func LevenshteinDistance(s, t string) int {
	m := len(s)
	n := len(t)

	// Handle edge cases
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Create the distance matrix with correct dimensions
	width := n + 1
	d := make([]int, (m+1)*width)

	// Initialize first row and column
	for i := 0; i <= m; i++ {
		d[i*width+0] = i
	}
	for j := 0; j <= n; j++ {
		d[0*width+j] = j
	}

	// Fill the dynamic programming table
	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			if s[i-1] == t[j-1] {
				d[i*width+j] = d[(i-1)*width+(j-1)]
			} else {
				d[i*width+j] = Min(
					d[(i-1)*width+j]+1,     // deletion
					d[i*width+(j-1)]+1,     // insertion
					d[(i-1)*width+(j-1)]+1) // substitution
			}
		}
	}

	return d[m*width+n]
}

// Min returns the minimum value from a slice of integers.
func Min(a ...int) int {
	if len(a) == 0 {
		return 0
	}
	min := a[0]
	for _, i := range a[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

// Max returns the maximum value from a slice of integers.
func Max(a ...int) int {
	if len(a) == 0 {
		return 0
	}
	max := a[0]
	for _, i := range a[1:] {
		if i > max {
			max = i
		}
	}
	return max
}

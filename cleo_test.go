package cleo

import "testing"

func TestLevenshtein(t *testing.T) {
	if LevenshteinDistance("abcdefghij", "abcdefghix") != 1 {
		t.Fail()
	}

	if LevenshteinDistance("abcdefghij", "abcdefghijk") != 1 {
		t.Fail()
	}

	if LevenshteinDistance("abcdefghij", "abcdefghi") != 1 {
		t.Fail()
	}
}

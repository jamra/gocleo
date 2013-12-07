package cleo

import "testing"

func TestScore(t *testing.T) {
	if Score("abcdefghij", "abcdefghix") != 0.90 {
		t.Fail()
	}
}

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

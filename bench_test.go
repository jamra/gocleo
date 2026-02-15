package cleo

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkSimpleCorpusLoading(b *testing.B) {
	// Create a simple test file
	filename := "test_corpus.txt"
	file, _ := os.Create(filename)
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(file, "word%03d\n", i)
	}
	file.Close()
	defer os.Remove(filename)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildIndexes(filename, nil)
	}
}

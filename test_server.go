package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jamra/gocleo/internal/fst"
)

type SearchResponse struct {
	Results []string `json:"results"`
	Count   int      `json:"count"`
	Time    string   `json:"time"`
}

func main() {
	// Create FST with sample data
	builder := fst.NewFSTBuilder()
	sampleData := []string{
		"apple", "banana", "cherry", "date", "elderberry",
		"fig", "grape", "honeydew", "kiwi", "lemon",
	}

	for _, item := range sampleData {
		builder.Add([]byte(item), uint64(len(item)))
	}

	searchFST, err := builder.Build()
	if err != nil {
		log.Fatal("Failed to build FST:", err)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		query := r.URL.Query().Get("q")
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if query == "" {
			json.NewEncoder(w).Encode(SearchResponse{
				Results: []string{},
				Count:   0,
				Time:    time.Since(start).String(),
			})
			return
		}

		// Prefix search
		var results []string
		iter := searchFST.PrefixIterator([]byte(query))
		for iter.HasNext() {
			key, _ := iter.Next()
			results = append(results, string(key))
		}

		response := SearchResponse{
			Results: results,
			Count:   len(results),
			Time:    time.Since(start).String(),
		}

		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<h1>FST Search Demo - Port 8081</h1>
<p>Try: <a href="/search?q=app">/search?q=app</a></p>
<p>Try: <a href="/search?q=ban">/search?q=ban</a></p>`)
	})

	fmt.Println("🚀 FST HTTP Server starting on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
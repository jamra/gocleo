package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	builder := fst.NewFSTBuilder()
	words := []string{"apple", "application", "apply", "go", "golang", "github", "test", "search", "server", "simple"}
	for i, word := range words {
		builder.Insert(word, uint64(i))
	}
	fstInstance, _ := builder.Finish()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html><head><title>FST Search Demo</title>
<style>body{font-family:Arial;max-width:600px;margin:50px auto;padding:20px}
input{padding:10px;font-size:16px;width:200px;margin:10px}
button{padding:10px 20px;font-size:16px;background:#007cba;color:white;border:none;cursor:pointer}
.result{padding:8px;margin:4px 0;background:#f0f8ff;border-radius:4px}</style>
</head><body>
<h1>🚀 FST Search Demo (Port 8081)</h1>
<div>
<input id="query" placeholder="Try: app, go, test..." value="app">
<button onclick="search()">Search</button>
</div>
<div id="results"></div>
<p><strong>API:</strong> <a href="/search?q=app">/search?q=app</a></p>
<script>
function search() {
	const q = document.getElementById("query").value;
	fetch("/search?q=" + q)
	.then(r => r.json())
	.then(d => {
		const html = "<h3>Results: " + d.count + "</h3>" + 
			d.results.map(r => "<div class=\"result\">" + r + "</div>").join("");
		document.getElementById("results").innerHTML = html;
	});
}
document.getElementById("query").addEventListener("keypress", function(e) {
	if (e.key === "Enter") search();
});
search(); // Initial search
</script></body></html>`
		w.Write([]byte(html))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		var results []string
		if query != "" {
			iter := fstInstance.Search(fst.StartsWith(query))
			for iter.Next() {
				results = append(results, iter.Key())
				if len(results) >= 10 { break }
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"query": query, "results": results, "count": len(results),
		})
		fmt.Printf("Search: %s -> %d results
", query, len(results))
	})

	fmt.Println("🚀 FST HTTP Server running on http://localhost:8081")
	fmt.Println("   Try: http://localhost:8081/")
	fmt.Println("   API: http://localhost:8081/search?q=your_query")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

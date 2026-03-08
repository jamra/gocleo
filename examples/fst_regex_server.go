package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jamra/gocleo/internal/fst"
	"github.com/jamra/gocleo/internal/scoring"
)

type SearchResult struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
	Timing  string         `json:"timing"`
	Error   string         `json:"error,omitempty"`
}

func main() {
	fmt.Println("Starting FST Regex Server on port 8081...")

	// Sample data for demonstration
	documents := []string{
		"apple fruit red sweet",
		"banana yellow tropical",
		"cherry small red fruit", 
		"date palm fruit",
		"elderberry dark purple",
		"fig Mediterranean fruit",
		"grape wine purple green",
		"honey sweet golden",
		"ice frozen water",
		"lemon citrus yellow sour",
		"application development",
		"programming tutorial",
		"testing framework",
		"debugging techniques",
		"optimization performance",
		"implementation details",
		"baking cookies delicious",
		"making bread fresh",
		"cooking pasta italian",
		"brewing coffee morning",
	}

	// Build FST
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		log.Fatalf("Failed to build FST: %v", err)
	}

	searchEngine := fst.NewSearchEngine(fstIndex, documents, scoring.DefaultScore)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head>
    <title>FST Regex Search Demo</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1000px; margin: 0 auto; padding: 20px; background: #f9f9f9; }
        .container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .search-section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; background: #fafafa; }
        .search-box { margin: 10px 0; }
        input { padding: 10px; width: 300px; font-size: 16px; margin-right: 10px; border: 1px solid #ccc; border-radius: 4px; }
        button { padding: 10px 20px; font-size: 16px; background: #007cba; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a85; }
        .results { margin: 20px 0; }
        .result { padding: 10px; border-bottom: 1px solid #eee; background: white; margin: 5px 0; border-radius: 3px; }
        .stats { background: #e8f4fd; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #007cba; }
        .examples { background: #f0f8ff; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .examples h4 { margin-top: 0; color: #333; }
        .example-button { background: #28a745; color: white; border: none; padding: 5px 10px; margin: 2px; cursor: pointer; border-radius: 3px; font-size: 12px; }
        .example-button:hover { background: #218838; }
        .section-title { color: #007cba; border-bottom: 2px solid #007cba; padding-bottom: 5px; }
        .api-info { background: #fff3cd; padding: 10px; margin: 10px 0; border-radius: 5px; border: 1px solid #ffeaa7; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="section-title">🔍 FST Regex Search Demo</h1>
        
        <div class="api-info">
            <strong>API Endpoints:</strong>
            <br>/search?q= - Basic search
            <br>/search/regex?pattern= - Regex search  
            <br>/search/complex - Complex search with multiple options
        </div>
        
        <div class="search-section">
            <h3>📝 Basic Search</h3>
            <div class="search-box">
                <input type="text" id="query" placeholder="Enter search term (e.g., 'app', 'fruit')..." onkeyup="if(event.key==='Enter')search()">
                <button onclick="search()">Search</button>
            </div>
        </div>
        
        <div class="search-section">
            <h3>🎯 Regex Search</h3>
            <div class="search-box">
                <input type="text" id="regex" placeholder="Enter regex pattern..." onkeyup="if(event.key==='Enter')regexSearch()">
                <button onclick="regexSearch()">Regex Search</button>
            </div>
            <div class="examples">
                <h4>Example Patterns:</h4>
                <button class="example-button" onclick="setRegex('app.*')" title="Words starting with 'app'">app.*</button>
                <button class="example-button" onclick="setRegex('.*ing$')" title="Words ending with 'ing'">.*ing$</button>
                <button class="example-button" onclick="setRegex('^[a-c].*')" title="Words starting with a, b, or c">^[a-c].*</button>
                <button class="example-button" onclick="setRegex('.*fruit.*')" title="Words containing 'fruit'">.*fruit.*</button>
                <button class="example-button" onclick="setRegex('^.{4}$')" title="Exactly 4 characters">^.{4}$</button>
                <button class="example-button" onclick="setRegex('.*[aeiou]{2}.*')" title="Double vowels">.*[aeiou]{2}.*</button>
                <button class="example-button" onclick="setRegex('^[^aeiou].*')" title="Start with consonant">^[^aeiou].*</button>
            </div>
        </div>
        
        <div class="search-section">
            <h3>🚀 Complex Search</h3>
            <div class="search-box">
                <input type="text" id="prefix" placeholder="Prefix..." style="width: 150px;">
                <input type="text" id="regexPattern" placeholder="Regex pattern..." style="width: 150px;">
                <input type="text" id="fuzzy" placeholder="Fuzzy term..." style="width: 120px;">
                <input type="number" id="distance" placeholder="Distance" style="width: 80px;" value="1">
                <button onclick="complexSearch()">Complex Search</button>
            </div>
            <div class="examples">
                <h4>Try these combinations:</h4>
                <button class="example-button" onclick="setComplex('app', '.*ion$', '', '1')" title="app + ending with ion">app + .*ion$</button>
                <button class="example-button" onclick="setComplex('', '.*ing$', 'making', '2')" title="ending with ing + fuzzy making">.*ing$ + fuzzy:making</button>
                <button class="example-button" onclick="setComplex('test', '.*', '', '0')" title="prefix test">prefix:test</button>
            </div>
        </div>
        
        <div id="results" class="results"></div>
        <div id="stats" class="stats" style="display: none;"></div>
    </div>
    
    <script>
        function search() {
            const query = document.getElementById('query').value;
            if (!query) return;
            
            fetch('/search?q=' + encodeURIComponent(query))
                .then(response => response.json())
                .then(showResults)
                .catch(error => console.error('Error:', error));
        }
        
        function regexSearch() {
            const regex = document.getElementById('regex').value;
            if (!regex) return;
            
            fetch('/search/regex?pattern=' + encodeURIComponent(regex))
                .then(response => response.json())
                .then(showResults)
                .catch(error => console.error('Error:', error));
        }
        
        function complexSearch() {
            const prefix = document.getElementById('prefix').value;
            const regexPattern = document.getElementById('regexPattern').value;
            const fuzzy = document.getElementById('fuzzy').value;
            const distance = document.getElementById('distance').value || 1;
            
            let url = '/search/complex?';
            if (prefix) url += 'prefix=' + encodeURIComponent(prefix) + '&';
            if (regexPattern) url += 'regex=' + encodeURIComponent(regexPattern) + '&';
            if (fuzzy) url += 'fuzzy=' + encodeURIComponent(fuzzy) + '&distance=' + distance + '&';
            
            fetch(url)
                .then(response => response.json())
                .then(showResults)
                .catch(error => console.error('Error:', error));
        }
        
        function setRegex(pattern) {
            document.getElementById('regex').value = pattern;
            regexSearch();
        }
        
        function setComplex(prefix, regex, fuzzy, distance) {
            document.getElementById('prefix').value = prefix;
            document.getElementById('regexPattern').value = regex;
            document.getElementById('fuzzy').value = fuzzy;
            document.getElementById('distance').value = distance;
            complexSearch();
        }
        
        function showResults(data) {
            let html = '<h3>🎯 Results (' + data.count + '):</h3>';
            if (data.error) {
                html += '<div style="color: red; padding: 10px; background: #ffe6e6; border-radius: 5px;">❌ Error: ' + data.error + '</div>';
            } else if (data.count === 0) {
                html += '<div style="color: #666; padding: 10px; background: #f5f5f5; border-radius: 5px;">No matches found</div>';
            } else {
                data.results.forEach((result, index) => {
                    html += '<div class="result">' + (index + 1) + '. ' + result.word + ' <span style="color: #666;">(score: ' + result.score.toFixed(2) + ')</span></div>';
                });
            }
            document.getElementById('results').innerHTML = html;
            
            let statsHtml = '<strong>⚡ Search Stats:</strong> Found ' + data.count + ' results in ' + data.timing;
            document.getElementById('stats').innerHTML = statsHtml;
            document.getElementById('stats').style.display = 'block';
        }
        
        // Load initial example
        setTimeout(() => {
            document.getElementById('query').value = 'app';
            search();
        }, 500);
    </script>
</body>
</html>`
		fmt.Fprint(w, html)
	})

	// Basic search endpoint
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		query := r.URL.Query().Get("q")
		start := time.Now()
		
		results := searchEngine.Search(query)
		duration := time.Since(start)
		
		response := SearchResponse{
			Results: make([]SearchResult, len(results)),
			Count:   len(results),
			Timing:  duration.String(),
		}
		
		for i, result := range results {
			response.Results[i] = SearchResult{
				Word:  result.Word,
				Score: result.Score,
			}
		}
		
		json.NewEncoder(w).Encode(response)
	})

	// Regex search endpoint
	http.HandleFunc("/search/regex", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		pattern := r.URL.Query().Get("pattern")
		start := time.Now()
		
		results, err := searchEngine.RegexSearch(pattern)
		duration := time.Since(start)
		
		response := SearchResponse{
			Results: make([]SearchResult, len(results)),
			Count:   len(results),
			Timing:  duration.String(),
		}
		
		if err != nil {
			response.Error = err.Error()
		} else {
			for i, result := range results {
				response.Results[i] = SearchResult{
					Word:  result.Word,
					Score: result.Score,
				}
			}
		}
		
		json.NewEncoder(w).Encode(response)
	})

	// Complex search endpoint
	http.HandleFunc("/search/complex", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		params := r.URL.Query()
		options := fst.ComplexSearchOptions{
			Prefix:       params.Get("prefix"),
			RegexPattern: params.Get("regex"),
			FuzzyPattern: params.Get("fuzzy"),
			Limit:        10, // Default limit
		}
		
		if distStr := params.Get("distance"); distStr != "" {
			if dist, err := strconv.Atoi(distStr); err == nil {
				options.FuzzyMaxDistance = dist
			}
		}
		
		start := time.Now()
		results, err := searchEngine.ComplexSearch(options)
		duration := time.Since(start)
		
		response := SearchResponse{
			Results: make([]SearchResult, len(results)),
			Count:   len(results),
			Timing:  duration.String(),
		}
		
		if err != nil {
			response.Error = err.Error()
		} else {
			for i, result := range results {
				response.Results[i] = SearchResult{
					Word:  result.Word,
					Score: result.Score,
				}
			}
		}
		
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("🚀 FST Regex Server ready!")
	fmt.Println("📍 URL: http://localhost:8081")
	fmt.Println("🎯 Features: Basic search, Regex patterns, Complex queries, Fuzzy matching")
	fmt.Println("")
	fmt.Println("📚 Example API calls:")
	fmt.Println("  http://localhost:8081/search?q=app")
	fmt.Println("  http://localhost:8081/search/regex?pattern=.*ing$")
	fmt.Println("  http://localhost:8081/search/complex?prefix=app&regex=.*ion$")
	fmt.Println("")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
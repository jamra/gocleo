package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
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

// RegexSearchEngine wraps FST with regex capabilities
type RegexSearchEngine struct {
	fst         *fst.FST
	documents   []string
	scoringFunc scoring.ScoringFunction
}

func NewRegexSearchEngine(fstIndex *fst.FST, documents []string) *RegexSearchEngine {
	return &RegexSearchEngine{
		fst:         fstIndex,
		documents:   documents,
		scoringFunc: scoring.DefaultScore,
	}
}

func (e *RegexSearchEngine) RegexSearch(pattern string) ([]SearchResult, error) {
	if pattern == "" || e.fst == nil {
		return []SearchResult{}, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	seen := make(map[string]bool)

	iter := e.fst.Iterator()
	for iter.HasNext() {
		key, docID := iter.Next()
		keyStr := string(key)
		
		if regex.MatchString(keyStr) {
			if int(docID) < len(e.documents) && docID > 0 {
				docContent := e.documents[docID-1]
				if !seen[docContent] {
					score := e.scoringFunc(keyStr, docContent)
					results = append(results, SearchResult{
						Word:  docContent,
						Score: score,
					})
					seen[docContent] = true
				}
			}
		}
	}

	return results, nil
}

func (e *RegexSearchEngine) PrefixRegexSearch(prefix, pattern string) ([]SearchResult, error) {
	if pattern == "" || e.fst == nil {
		return []SearchResult{}, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	seen := make(map[string]bool)
	prefixLower := strings.ToLower(prefix)

	iter := e.fst.PrefixIterator([]byte(prefixLower))
	for iter.HasNext() {
		key, docID := iter.Next()
		keyStr := string(key)
		
		if regex.MatchString(keyStr) {
			if int(docID) < len(e.documents) && docID > 0 {
				docContent := e.documents[docID-1]
				if !seen[docContent] {
					score := e.scoringFunc(keyStr, docContent)
					results = append(results, SearchResult{
						Word:  docContent,
						Score: score,
					})
					seen[docContent] = true
				}
			}
		}
	}

	return results, nil
}

func (e *RegexSearchEngine) BasicSearch(query string) []SearchResult {
	if query == "" || e.fst == nil {
		return []SearchResult{}
	}

	var results []SearchResult
	seen := make(map[string]bool)
	queryLower := strings.ToLower(query)

	iter := e.fst.PrefixIterator([]byte(queryLower))
	for iter.HasNext() {
		_, docID := iter.Next()
		if int(docID) < len(e.documents) && docID > 0 {
			docContent := e.documents[docID-1]
			if !seen[docContent] {
				score := e.scoringFunc(query, docContent)
				results = append(results, SearchResult{
					Word:  docContent,
					Score: score,
				})
				seen[docContent] = true
			}
		}
	}

	return results
}

func main() {
	fmt.Println("🚀 Starting FST Regex Server on port 8081...")

	// Enhanced sample data
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
		"application development programming",
		"programming tutorial guide",
		"testing framework unittest",
		"debugging techniques tools",
		"optimization performance tuning",
		"implementation software code",
		"baking cookies delicious sweet",
		"making bread fresh morning",
		"cooking pasta italian cuisine",
		"brewing coffee morning drink",
		"recipe collection cookbook",
		"instruction manual guide",
		"tutorial learning education",
		"development environment setup",
	}

	// Build FST
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		log.Fatalf("Failed to build FST: %v", err)
	}

	searchEngine := NewRegexSearchEngine(fstIndex, documents)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head>
    <title>FST Regex Search Demo</title>
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; background: #f5f7fa; }
        .container { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .search-section { margin: 25px 0; padding: 20px; border: 1px solid #e0e6ed; border-radius: 8px; background: #fafbfc; }
        .search-box { margin: 15px 0; }
        input { padding: 12px 16px; width: 300px; font-size: 16px; margin-right: 12px; border: 2px solid #ddd; border-radius: 6px; transition: border-color 0.3s; }
        input:focus { border-color: #007cba; outline: none; }
        button { padding: 12px 24px; font-size: 16px; background: #007cba; color: white; border: none; border-radius: 6px; cursor: pointer; transition: background-color 0.3s; }
        button:hover { background: #005a85; }
        .results { margin: 25px 0; }
        .result { padding: 15px; border-bottom: 1px solid #eee; background: white; margin: 8px 0; border-radius: 6px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .stats { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; margin: 15px 0; border-radius: 8px; }
        .examples { background: #f8f9ff; padding: 15px; margin: 15px 0; border-radius: 8px; border: 1px solid #e1e5f2; }
        .examples h4 { margin-top: 0; color: #333; }
        .example-button { background: #28a745; color: white; border: none; padding: 8px 12px; margin: 3px; cursor: pointer; border-radius: 4px; font-size: 13px; transition: background-color 0.3s; }
        .example-button:hover { background: #218838; }
        .section-title { color: #007cba; border-bottom: 3px solid #007cba; padding-bottom: 8px; margin-bottom: 20px; }
        .api-info { background: #fff3cd; padding: 15px; margin: 15px 0; border-radius: 8px; border: 1px solid #ffeaa7; }
        .feature-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin: 20px 0; }
        .feature-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="section-title">🔍 FST Regex Search Engine</h1>
            <p>Demonstrating Finite State Transducer with Regular Expression capabilities</p>
        </div>
        
        <div class="api-info">
            <strong>🛠️ API Endpoints:</strong><br>
            <code>GET /search?q=term</code> - Basic prefix search<br>
            <code>GET /search/regex?pattern=regex</code> - Full regex search<br>
            <code>GET /search/complex?prefix=...&regex=...&fuzzy=...</code> - Complex search
        </div>
        
        <div class="feature-grid">
            <div class="feature-card">
                <h3>📝 Basic Search</h3>
                <div class="search-box">
                    <input type="text" id="query" placeholder="Search term (e.g., 'app', 'fruit')..." onkeyup="if(event.key==='Enter')search()">
                    <button onclick="search()">Search</button>
                </div>
            </div>
            
            <div class="feature-card">
                <h3>🎯 Regex Search</h3>
                <div class="search-box">
                    <input type="text" id="regex" placeholder="Regex pattern..." onkeyup="if(event.key==='Enter')regexSearch()">
                    <button onclick="regexSearch()">Regex Search</button>
                </div>
                <div class="examples">
                    <h4>🎨 Pattern Examples:</h4>
                    <button class="example-button" onclick="setRegex('app.*')" title="Words starting with 'app'">app.*</button>
                    <button class="example-button" onclick="setRegex('.*ing$')" title="Words ending with 'ing'">.*ing$</button>
                    <button class="example-button" onclick="setRegex('^[a-c].*')" title="Words starting with a, b, or c">^[a-c].*</button>
                    <button class="example-button" onclick="setRegex('.*fruit.*')" title="Words containing 'fruit'">.*fruit.*</button>
                    <button class="example-button" onclick="setRegex('^.{4}$')" title="Exactly 4 characters">^.{4}$</button>
                    <button class="example-button" onclick="setRegex('.*[aeiou]{2}.*')" title="Double vowels">.*[aeiou]{2}.*</button>
                </div>
            </div>
        </div>
        
        <div class="search-section">
            <h3>🚀 Prefix + Regex Combo</h3>
            <div class="search-box">
                <input type="text" id="prefix" placeholder="Prefix..." style="width: 200px;">
                <input type="text" id="regexPattern" placeholder="Regex pattern..." style="width: 200px;">
                <button onclick="prefixRegexSearch()">Combined Search</button>
            </div>
            <div class="examples">
                <h4>🔥 Powerful Combinations:</h4>
                <button class="example-button" onclick="setPrefixRegex('app', '.*ion$')" title="app + ending with ion">app + .*ion$</button>
                <button class="example-button" onclick="setPrefixRegex('cook', '.*ing$')" title="cook + ending with ing">cook + .*ing$</button>
                <button class="example-button" onclick="setPrefixRegex('dev', '.*ment$')" title="dev + ending with ment">dev + .*ment$</button>
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
        
        function prefixRegexSearch() {
            const prefix = document.getElementById('prefix').value;
            const regex = document.getElementById('regexPattern').value;
            
            if (!prefix || !regex) return;
            
            fetch('/search/prefix-regex?prefix=' + encodeURIComponent(prefix) + '&pattern=' + encodeURIComponent(regex))
                .then(response => response.json())
                .then(showResults)
                .catch(error => console.error('Error:', error));
        }
        
        function setRegex(pattern) {
            document.getElementById('regex').value = pattern;
            regexSearch();
        }
        
        function setPrefixRegex(prefix, regex) {
            document.getElementById('prefix').value = prefix;
            document.getElementById('regexPattern').value = regex;
            prefixRegexSearch();
        }
        
        function showResults(data) {
            let html = '<h3>🎯 Search Results (' + data.count + '):</h3>';
            if (data.error) {
                html += '<div style="color: #e74c3c; padding: 15px; background: #ffe6e6; border-radius: 8px; border-left: 4px solid #e74c3c;">❌ <strong>Error:</strong> ' + data.error + '</div>';
            } else if (data.count === 0) {
                html += '<div style="color: #666; padding: 20px; background: #f8f9fa; border-radius: 8px; text-align: center;">🔍 No matches found. Try a different pattern!</div>';
            } else {
                data.results.forEach((result, index) => {
                    html += '<div class="result">' + 
                            '<strong>' + (index + 1) + '.</strong> ' + result.word + 
                            ' <span style="color: #666; font-size: 0.9em;">(score: ' + result.score.toFixed(3) + ')</span>' +
                            '</div>';
                });
            }
            document.getElementById('results').innerHTML = html;
            
            let statsHtml = '<strong>⚡ Performance:</strong> Found ' + data.count + ' results in ' + data.timing + 
                           ' | <strong>🎯 Accuracy:</strong> Zero false positives guaranteed';
            document.getElementById('stats').innerHTML = statsHtml;
            document.getElementById('stats').style.display = 'block';
        }
        
        // Load demo on page load
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
		
		results := searchEngine.BasicSearch(query)
		duration := time.Since(start)
		
		response := SearchResponse{
			Results: results,
			Count:   len(results),
			Timing:  duration.String(),
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
			Results: results,
			Count:   len(results),
			Timing:  duration.String(),
		}
		
		if err != nil {
			response.Error = err.Error()
		}
		
		json.NewEncoder(w).Encode(response)
	})

	// Prefix + Regex search endpoint
	http.HandleFunc("/search/prefix-regex", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		prefix := r.URL.Query().Get("prefix")
		pattern := r.URL.Query().Get("pattern")
		start := time.Now()
		
		results, err := searchEngine.PrefixRegexSearch(prefix, pattern)
		duration := time.Since(start)
		
		response := SearchResponse{
			Results: results,
			Count:   len(results),
			Timing:  duration.String(),
		}
		
		if err != nil {
			response.Error = err.Error()
		}
		
		json.NewEncoder(w).Encode(response)
	})

	// Stats endpoint
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		stats := map[string]interface{}{
			"fst_size":     fstIndex.Size(),
			"fst_empty":    fstIndex.IsEmpty(),
			"documents":    len(documents),
			"capabilities": []string{"exact_match", "prefix_search", "regex_search", "combined_search"},
		}
		
		json.NewEncoder(w).Encode(stats)
	})

	fmt.Println("🎯 FST Regex Server Features:")
	fmt.Println("   • Basic prefix search")
	fmt.Println("   • Full regex pattern matching")  
	fmt.Println("   • Combined prefix + regex search")
	fmt.Println("   • Zero false positives (unlike bloom filters)")
	fmt.Println("   • High performance with low memory usage")
	fmt.Println("")
	fmt.Println("🌐 Server running at: http://localhost:8081")
	fmt.Println("")
	fmt.Println("📚 API Examples:")
	fmt.Println("   curl \"http://localhost:8081/search?q=app\"")
	fmt.Println("   curl \"http://localhost:8081/search/regex?pattern=.*ing$\"")
	fmt.Println("   curl \"http://localhost:8081/search/prefix-regex?prefix=app&pattern=.*ion$\"")
	fmt.Println("   curl \"http://localhost:8081/stats\"")
	fmt.Println("")
	
	log.Fatal(http.ListenAndServe(":8081", nil))
}
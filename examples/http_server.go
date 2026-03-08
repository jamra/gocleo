package main

import (
	"log"
	"os"

	"github.com/jamra/gocleo/api/http"
	"github.com/jamra/gocleo/pkg/cleo"
)

func main() {
	// Create a sample data directory and index
	tempDir := "/tmp/cleo_http_example"
	os.RemoveAll(tempDir) // Clean up any existing data
	os.MkdirAll(tempDir, 0755)

	// Initialize Cleo client
	client, err := cleo.NewClient(tempDir)
	if err != nil {
		log.Fatalf("Failed to create Cleo client: %v", err)
	}
	defer client.Close()

	// Add sample documents
	sampleDocs := []struct {
		id      string
		content string
	}{
		{"doc1", "apple fruit red delicious healthy snack"},
		{"doc2", "banana yellow fruit tropical potassium"},
		{"doc3", "orange citrus vitamin c juice breakfast"},
		{"doc4", "pizza italian food cheese tomato sauce"},
		{"doc5", "burger american fast food beef lettuce"},
		{"doc6", "sushi japanese raw fish rice seaweeeed"},
		{"doc7", "pasta italian noodles sauce garlic"},
		{"doc8", "salad healthy greens vegetables dressing"},
		{"doc9", "chocolate sweet dessert cocoa milk"},
		{"doc10", "coffee morning drink caffeine energy"},
	}

	log.Println("Adding sample documents to index...")
	for _, doc := range sampleDocs {
		if err := client.Index(doc.id, doc.content); err != nil {
			log.Printf("Failed to index document %s: %v", doc.id, err)
		}
	}

	log.Println("🚀 Starting HTTP server on port 8081 (avoiding conflict with port 8080)")
	log.Println("")
	log.Println("Try these URLs:")
	log.Println("  http://localhost:8081/search?q=fruit")
	log.Println("  http://localhost:8081/search?q=italian") 
	log.Println("  http://localhost:8081/cleo/food")
	log.Println("  http://localhost:8081/stats")
	log.Println("")

	// Start HTTP server on port 8081 (not 8080 to avoid conflicts)
	if err := http.ListenAndServe("8081", client); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

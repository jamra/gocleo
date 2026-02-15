// Package main demonstrates the legacy/backward-compatible API.
// For new applications, consider using the new API in new_api_example.go
package main

import (
	"log"
	"github.com/jamra/gocleo"
)

func main() {
	// This uses the legacy API for backward compatibility
	log.Println("Starting Cleo search server using legacy API...")
	log.Println("This is maintained for backward compatibility.")
	log.Println("For new applications, see new_api_example.go")
	
	err := cleo.InitAndRun("./w1_fixed.txt", "9999", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

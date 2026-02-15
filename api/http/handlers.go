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

// Package http provides HTTP handlers for Cleo search functionality.
package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jamra/gocleo/pkg/cleo"
)

// Server wraps a Cleo client with HTTP server functionality.
type Server struct {
	client *cleo.Client
}

// NewServer creates a new HTTP server with the given Cleo client.
func NewServer(client *cleo.Client) *Server {
	return &Server{
		client: client,
	}
}

// SearchHandler handles search requests at /search?q=query or /search?query=query
func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for web applications
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		query = r.URL.Query().Get("query")
	}

	if query == "" {
		http.Error(w, `{"error": "Missing query parameter 'q' or 'query'"}`, http.StatusBadRequest)
		return
	}

	// Perform search
	results, err := s.client.Search(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Search failed: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Return results as JSON
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode results"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// StatsHandler provides statistics about the search index.
func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := s.client.GetStats()
	jsonResponse, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode stats"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// RegisterRoutes registers all HTTP routes on the given mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", s.SearchHandler)
	mux.HandleFunc("/stats", s.StatsHandler)
	
	// Legacy route for backward compatibility
	mux.HandleFunc("/cleo", s.LegacyCleoHandler)
}

// LegacyCleoHandler provides backward compatibility with the original /cleo endpoint.
func (s *Server) LegacyCleoHandler(w http.ResponseWriter, r *http.Request) {
	// Original handler expected /cleo?query=value or path like /cleo/query
	query := r.URL.Query().Get("query")
	
	if query == "" {
		// Try to extract from path (e.g., /cleo/pizza)
		path := r.URL.Path
		if len(path) > 6 && path[:6] == "/cleo/" {
			query = path[6:]
		}
	}

	if query == "" {
		http.Error(w, `{"error": "Missing query"}`, http.StatusBadRequest)
		return
	}

	// Perform search
	results, err := s.client.Search(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Search failed: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Return results as JSON (matches original format)
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode results"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// ListenAndServe starts an HTTP server on the specified port with Cleo search functionality.
func ListenAndServe(port string, client *cleo.Client) error {
	server := NewServer(client)
	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port: %s", port)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting Cleo search server on http://localhost%s", addr)
	log.Printf("Search endpoint: http://localhost%s/search?q=your_query", addr)
	log.Printf("Legacy endpoint: http://localhost%s/cleo/your_query", addr)
	log.Printf("Stats endpoint: http://localhost%s/stats", addr)

	return http.ListenAndServe(addr, mux)
}

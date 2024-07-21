package handlers

import (
    "encoding/json"
    "net/http"
	"log"

	"github.com/tmohagan/go-search-service/db"
)

type SearchResult struct {
    ID      string `json:"id" bson:"_id"`
    Title   string `json:"title"`
    Summary string `json:"summary"`
    Type    string `json:"type"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Missing search query", http.StatusBadRequest)
        return
    }
    
    log.Printf("Received search query: %s", query)

    results, err := db.PerformSearch(query)
    if err != nil {
        log.Printf("Error performing search: %v", err)
        http.Error(w, "Error performing search", http.StatusInternalServerError)
        return
    }
    
    log.Printf("Found %d results", len(results))

    // Convert results to SearchResult structs
    var searchResults []SearchResult
    for _, result := range results {
        searchResult := SearchResult{
            ID:      result["_id"].(string),
            Title:   result["title"].(string),
            Summary: result["summary"].(string),
            Type:    result["type"].(string),
        }
        searchResults = append(searchResults, searchResult)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(searchResults)
}
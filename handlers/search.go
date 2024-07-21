package handlers

import (
    "encoding/json"
    "net/http"
    "log"

    "github.com/tmohagan/go-search-service/db"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchResult struct {
    ID      primitive.ObjectID `json:"id" bson:"_id"`
    Title   string             `json:"title"`
    Summary string             `json:"summary"`
    Type    string             `json:"type"`
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
        var searchResult SearchResult
        
        // Handle _id as ObjectID
        if oid, ok := result["_id"].(primitive.ObjectID); ok {
            searchResult.ID = oid
        } else {
            log.Printf("Unexpected type for _id: %T", result["_id"])
            continue
        }

        // Handle other fields
        if title, ok := result["title"].(string); ok {
            searchResult.Title = title
        }
        if summary, ok := result["summary"].(string); ok {
            searchResult.Summary = summary
        }
        if typeStr, ok := result["type"].(string); ok {
            searchResult.Type = typeStr
        }

        searchResults = append(searchResults, searchResult)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(searchResults)
}
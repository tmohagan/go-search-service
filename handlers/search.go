package handlers

import (
    "encoding/json"
    "net/http"
    "log"
    "strconv"
    "strings"

    "github.com/tmohagan/go-search-service/db"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchResult struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Summary string `json:"summary"`
    Type    string `json:"type"`
}

type SearchResponse struct {
    Results     []SearchResult `json:"results"`
    TotalCount  int64          `json:"totalCount"`
    CurrentPage int64          `json:"currentPage"`
    TotalPages  int64          `json:"totalPages"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := sanitizeInput(r.URL.Query().Get("q"))
    if query == "" {
        http.Error(w, "Missing search query", http.StatusBadRequest)
        return
    }

    page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
    if err != nil || page < 1 {
        page = 1
    }
    
    resultsPerPage := int64(10)

    log.Printf("Received search query: %s, page: %d", query, page)

    results, totalCount, err := db.PerformSearch(query, page, resultsPerPage)
    if err != nil {
        log.Printf("Error performing search: %v", err)
        http.Error(w, "Error performing search", http.StatusInternalServerError)
        return
    }

    searchResults := make([]SearchResult, 0, len(results))
    for _, result := range results {
        id, ok := result["_id"].(primitive.ObjectID)
        if !ok {
            log.Printf("Unexpected type for _id: %T", result["_id"])
            continue
        }

        title, ok := result["title"].(string)
        if !ok {
            log.Printf("Unexpected type for title: %T", result["title"])
            continue
        }

        summary, ok := result["summary"].(string)
        if !ok {
            log.Printf("Unexpected type for summary: %T", result["summary"])
            continue
        }

        resultType, ok := result["type"].(string)
        if !ok {
            log.Printf("Unexpected type for type: %T", result["type"])
            continue
        }

        searchResults = append(searchResults, SearchResult{
            ID:      id.Hex(),
            Title:   title,
            Summary: summary,
            Type:    resultType,
        })
    }

    totalPages := (totalCount + resultsPerPage - 1) / resultsPerPage

    response := SearchResponse{
        Results:     searchResults,
        TotalCount:  totalCount,
        CurrentPage: page,
        TotalPages:  totalPages,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func sanitizeInput(input string) string {
    return strings.Map(func(r rune) rune {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' {
            return r
        }
        return -1
    }, input)
}
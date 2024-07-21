package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/tmohagan/go-search-service/db"
    "github.com/tmohagan/go-search-service/handlers"
)

func main() {
    err := db.ConnectDB()
    if err != nil {
        log.Fatal(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/search", handlers.SearchHandler).Methods("GET")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
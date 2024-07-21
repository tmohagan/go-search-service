package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
	"github.com/joho/godotenv"
    "github.com/tmohagan/go-search-service/db"
    "github.com/tmohagan/go-search-service/handlers"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    err = db.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
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
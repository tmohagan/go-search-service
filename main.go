package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/tmohagan/go-search-service/db"
    "github.com/tmohagan/go-search-service/handlers"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    if err := db.ConnectDB(); err != nil {
        return err
    }

    r := mux.NewRouter()
    r.Use(loggingMiddleware)
    r.HandleFunc("/search", handlers.SearchHandler).Methods("GET")

    port := getEnv("PORT", "8080")
    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Server starting on port %s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe(): %v", err)
        }
    }()

    return gracefulShutdown(srv)
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

func gracefulShutdown(srv *http.Server) error {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    srv.Shutdown(ctx)
    log.Println("shutting down")
    return nil
}
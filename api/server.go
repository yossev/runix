package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "runix/api/api"
)

func main() {
    r := mux.NewRouter()
    
    // Set up API routes
    api.SetupRoutes(r)
    
    // Start the server
    http.Handle("/", r)
    log.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

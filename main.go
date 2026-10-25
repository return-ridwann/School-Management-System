package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// sends a JSON response with status code
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("writeJSON error:", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "School Management System",
		})
	})

	// serve static files from directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

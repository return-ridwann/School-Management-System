package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sms/db"
)

// send JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("writeJSON error:", err)
	}
}

// check if server and db are healthy
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := db.DB.Ping(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"status": "error",
			"db":     "disconnected",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"db":     "connected",
	})
}

func main() {
	// connect to sqlite database
	if err := db.InitDB("school.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close()

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /api/health", handleHealth)

	// serve static files from directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

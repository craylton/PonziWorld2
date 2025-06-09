package main

import (
	"fmt"
	"net/http"
)

// HelloHandler handles the /api/hello endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello from Go backend!"}`)
}

// RegisterRoutes registers all API routes to the mux
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/hello", HelloHandler)
	// Add more routes here as your API grows
}

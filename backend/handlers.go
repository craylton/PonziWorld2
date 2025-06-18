package main

import (
	"net/http"
)

// RegisterRoutes registers all API routes to the mux
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user", CreateUserHandler)
	mux.HandleFunc("/api/login", LoginHandler)
}
